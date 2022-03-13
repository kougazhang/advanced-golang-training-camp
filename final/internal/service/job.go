package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/webmin7761/go-school/homework/final/api/fare/v1"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	"github.com/webmin7761/go-school/homework/final/internal/cache/redis"
	"github.com/webmin7761/go-school/homework/final/internal/mq/kafka"
)

func NewJobService(fare *biz.FareUsecase, cache *redis.RadixRC3, mq *kafka.KafkaClient, logger log.Logger) *JobService {
	return &JobService{
		ttl:   5000,
		fare:  fare,
		cache: cache,
		mq:    mq,
		log:   log.NewHelper(logger),
	}
}

func (j *JobService) UpdateCache(ctx context.Context) error {
	return j.mq.Consume(j.writeCache)
}

func (j *JobService) writeCache(ctx context.Context, msg string, err error) error {
	if err != nil {
		return err
	}

	var req pb.PricingRequest
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return err
	}

	p, err := j.fare.Pricing(ctx, &biz.Fare{
		OrgAirport:      req.OrgAirport,
		ArrAirport:      req.ArrAirport,
		PassageType:     req.PassageType.String(),
		FirstTravelDate: req.FlightDatetime.AsTime(),
		LastTravelDate:  req.FlightDatetime.AsTime(),
	})
	if err != nil {
		return err
	}

	key := genKey(req.OrgAirport, req.ArrAirport, req.FlightDatetime, req.PassageType)
	return j.cache.Set(key, strconv.FormatFloat(p.Amount, 'f', 2, 32), j.ttl)
}
