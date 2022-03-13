package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"github.com/webmin7761/go-school/homework/final/internal/conf"
	"github.com/webmin7761/go-school/homework/final/internal/service"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(service *service.JobService) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		service.UpdateCache(ctx)
		return nil
	}
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	config := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := config.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := config.Scan(&bc); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
	}(ctx)

	run, cleanup, err := initApp(bc.Data, bc.Cache, bc.Mq, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return run(ctx)
	})

	g.Go(func() error {
		return sigProcess()(ctx)
	})

	if err := g.Wait(); err != nil {
		logger.Log(log.LevelError, err)
	}
}

func sigProcess(sig ...os.Signal) func(context.Context) error {
	return func(ctx context.Context) error {

		if len(sig) == 0 {
			sig = append(sig, os.Interrupt)
		}

		done := make(chan os.Signal, len(sig))
		signal.Notify(done, sig...)

		var err error
		select {
		case <-ctx.Done():
		case s := <-done:
			err = errors.New("main: " + s.String())
		}

		signal.Stop(done)
		close(done)

		return err
	}
}
