// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	"github.com/webmin7761/go-school/homework/final/internal/cache"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
	"github.com/webmin7761/go-school/homework/final/internal/data"
	"github.com/webmin7761/go-school/homework/final/internal/mq"
	server "github.com/webmin7761/go-school/homework/final/internal/server/fare"
	"github.com/webmin7761/go-school/homework/final/internal/service"
	"go.opentelemetry.io/otel/trace"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *conf.Cache, *conf.MessageQueue, *conf.Service, trace.TracerProvider, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		cache.ProviderSet,
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		mq.ProviderSet,
		newApp))
}
