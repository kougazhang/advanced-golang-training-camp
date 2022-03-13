// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
	travel "github.com/webmin7761/go-school/homework/final/internal/server/travel"
	"github.com/webmin7761/go-school/homework/final/internal/service"
	"go.opentelemetry.io/otel/trace"
)

// initApp init kratos application.
func initApp(*conf.Server, trace.TracerProvider, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		travel.ProviderSet,
		service.TravelProviderSet,
		newApp))
}
