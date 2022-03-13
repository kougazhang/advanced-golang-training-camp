package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	f "github.com/webmin7761/go-school/homework/final/api/fare/v1"
	s "github.com/webmin7761/go-school/homework/final/api/shop/v1"
	t "github.com/webmin7761/go-school/homework/final/api/travel/v1"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	c "github.com/webmin7761/go-school/homework/final/internal/cache"
	m "github.com/webmin7761/go-school/homework/final/internal/mq"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewFareService)

var ShopProviderSet = wire.NewSet(NewShoppingService)

var TravelProviderSet = wire.NewSet(NewTravelService)

var JobProviderSet = wire.NewSet(NewJobService)

type FareService struct {
	f.UnimplementedFareServiceServer

	log   *log.Helper
	fare  *biz.FareUsecase
	cache c.Cache
	mq    m.MessageQueue
}

type ShoppingService struct {
	s.UnimplementedShopServiceServer

	log    *log.Helper
	fare   f.FareServiceHTTPClient
	travel t.TravelServiceHTTPClient
}

type TravelService struct {
	t.UnimplementedTravelServiceServer
	log *log.Helper
}

type JobService struct {
	ttl   int
	fare  *biz.FareUsecase
	log   *log.Helper
	cache c.Cache
	mq    m.MessageQueue
}
