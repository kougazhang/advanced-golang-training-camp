package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	pb "github.com/webmin7761/go-school/homework/final/api/fare/v1"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
)

func NewFareClient(conf *conf.Service) pb.FareServiceHTTPClient {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint(conf.ServiceMap["Fare"]),
	)
	if err != nil {
		panic(err)
	}

	return pb.NewFareServiceHTTPClient(conn)
}
