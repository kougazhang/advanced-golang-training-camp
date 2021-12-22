package homework

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type GroupServer struct {
	servers []*Server
	g       errgroup.Group
	quit    chan struct{}
	wg      sync.WaitGroup
}

type Server struct {
	Name   string
	Server *http.Server
	// https 服务需配置以下两项
	CertFile string
	KeyFile  string
}

func (srv *Server) start() error {
	log.Printf("server %s starting ...\n", srv.Name)
	if len(srv.CertFile) > 0 && len(srv.KeyFile) > 0 {
		return srv.Server.ListenAndServeTLS(srv.CertFile, srv.KeyFile)
	}
	return srv.Server.ListenAndServe()
}

func (srv *Server) shutdown() error {
	log.Printf("server %s shutdowning ...\n", srv.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Server.Shutdown(ctx)
}

// Register: 注册服务
func (s *GroupServer) Register(server ...*Server) {
	s.servers = append(s.servers, server...)
}

// Run: 启动服务
func (s *GroupServer) Run() error {
	s.quit = make(chan struct{}, len(s.servers)*2)
	for _, srv := range s.servers {
		srv := srv
		s.wg.Add(2)
		// 启动 server
		s.g.Go(func() error {
			defer s.wg.Done()
			return srv.start()
		})
		// 收到信号关闭 server
		s.g.Go(func() error {
			defer s.wg.Done()
			<-s.quit
			return srv.shutdown()
		})
	}

	// 收到 linux Interrupt 信号时集体退出
	go func() {
		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt
		for i, max := 0, len(s.servers); i < max; i++ {
			s.quit <- struct{}{}
		}
	}()

	// 当某个 server 出错时集体退出
	go func() {
		if err := s.g.Wait(); err != nil {
			for i, max := 0, len(s.servers); i < max; i++ {
				s.quit <- struct{}{}
			}
		}
	}()

	s.wg.Wait()
	return nil
}
