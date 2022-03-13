package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	pb "github.com/webmin7761/go-school/homework/final/api/travel/v1"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
)

func NewTravelClient(conf *conf.Service) pb.TravelServiceHTTPClient {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint(conf.ServiceMap["Travel"]),
	)
	if err != nil {
		panic(err)
	}

	return pb.NewTravelServiceHTTPClient(conn)
}
