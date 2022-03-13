// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	"github.com/webmin7761/go-school/homework/final/internal/cache"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
	"github.com/webmin7761/go-school/homework/final/internal/data"
	"github.com/webmin7761/go-school/homework/final/internal/mq"
	"github.com/webmin7761/go-school/homework/final/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Data, *conf.Cache, *conf.MessageQueue, log.Logger) (func(context.Context) error, func(), error) {
	panic(wire.Build(
		cache.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.JobProviderSet,
		mq.ProviderSet,
		newApp))
}
