package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
	"github.com/webmin7761/go-school/homework/final/internal/data/ent"

	_ "github.com/mattn/go-sqlite3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewFareRepo)

// Data .
type Data struct {
	db *ent.Client
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	client, err := ent.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, nil, err
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}

	d := &Data{
		db: client,
	}
	return d, func() {
		log.Info("message", "closing the data resources")
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
