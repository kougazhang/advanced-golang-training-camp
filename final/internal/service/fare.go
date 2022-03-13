package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/webmin7761/go-school/homework/final/api/common/v1"
	pb "github.com/webmin7761/go-school/homework/final/api/fare/v1"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	"github.com/webmin7761/go-school/homework/final/internal/cache/redis"
	"github.com/webmin7761/go-school/homework/final/internal/mq/kafka"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func NewFareService(fare *biz.FareUsecase, cache *redis.RadixRC3, mq *kafka.KafkaClient, logger log.Logger) *FareService {
	return &FareService{
		fare:  fare,
		cache: cache,
		mq:    mq,
		log:   log.NewHelper(logger),
	}
}

func (s *FareService) CreateFare(ctx context.Context, req *pb.CreateFareRequest) (*pb.CreateFareReply, error) {
	s.log.Infof("input data %v", req)
	id, err := s.fare.Create(ctx, &biz.Fare{
		OrgAirport:      req.Fare.OrgAirport,
		ArrAirport:      req.Fare.ArrAirport,
		PassageType:     req.Fare.PassageType.String(),
		FirstTravelDate: req.Fare.FirstTravelDate.AsTime(),
		LastTravelDate:  req.Fare.LastTravelDate.AsTime(),
		Amount:          req.Fare.Amount.Value,
	})
	return &pb.CreateFareReply{
		Result: &common.Result{Code: "0"},
		Id:     id,
	}, err
}

func (s *FareService) UpdateFare(ctx context.Context, req *pb.UpdateFareRequest) (*pb.UpdateFareReply, error) {
	s.log.Infof("input data %v", req)
	err := s.fare.Update(ctx, req.Fare.Id, &biz.Fare{
		OrgAirport:      req.Fare.OrgAirport,
		ArrAirport:      req.Fare.ArrAirport,
		PassageType:     req.Fare.PassageType.String(),
		FirstTravelDate: req.Fare.FirstTravelDate.AsTime(),
		LastTravelDate:  req.Fare.LastTravelDate.AsTime(),
		Amount:          req.Fare.Amount.Value,
	})
	return &pb.UpdateFareReply{}, err
}

func (s *FareService) DeleteFare(ctx context.Context, req *pb.DeleteFareRequest) (*pb.DeleteFareReply, error) {
	s.log.Infof("input data %v", req)
	err := s.fare.Delete(ctx, req.Id)
	return &pb.DeleteFareReply{}, err
}

func (s *FareService) GetFare(ctx context.Context, req *pb.GetFareRequest) (*pb.GetFareReply, error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetFare")
	defer span.End()
	p, err := s.fare.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	psgType := common.PassageTypes(common.PassageTypes_value[p.PassageType])
	return &pb.GetFareReply{
		Result: &common.Result{Code: "0"},
		Fare: &pb.Fare{
			Id:              p.Id,
			OrgAirport:      p.OrgAirport,
			ArrAirport:      p.ArrAirport,
			FirstTravelDate: timestamppb.New(p.FirstTravelDate),
			LastTravelDate:  timestamppb.New(p.LastTravelDate),
			PassageType:     psgType,
			Amount:          wrapperspb.Double(p.Amount),
		},
	}, nil
}

func (s *FareService) Pricing(ctx context.Context, req *pb.PricingRequest) (*pb.PricingResponse, error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "PricingFare")
	defer span.End()

	key := genKey(req.OrgAirport, req.ArrAirport, req.FlightDatetime, req.PassageType)
	v, err := s.cache.Get(key)

	if v != "" && err == nil {
		amount, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, err
		}
		return &pb.PricingResponse{
			Result:         &common.Result{Code: "0"},
			OrgAirport:     req.OrgAirport,
			ArrAirport:     req.ArrAirport,
			FlightDatetime: req.FlightDatetime,
			PassageType:    req.PassageType,
			Amount:         wrapperspb.Double(amount),
		}, nil
	}

	p, err := s.fare.Pricing(ctx, &biz.Fare{
		OrgAirport:      req.OrgAirport,
		ArrAirport:      req.ArrAirport,
		PassageType:     req.PassageType.String(),
		FirstTravelDate: req.FlightDatetime.AsTime(),
		LastTravelDate:  req.FlightDatetime.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	res := &pb.PricingResponse{
		Result:         &common.Result{Code: "0"},
		OrgAirport:     req.OrgAirport,
		ArrAirport:     req.ArrAirport,
		FlightDatetime: req.FlightDatetime,
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(p.Amount),
	}

	msg, _ := json.Marshal(req)
	s.mq.Produce(string(msg))

	return res, nil
}

func (s *FareService) PriceCalendar(ctx context.Context, req *pb.PriceCalendarRequest) (*pb.PriceCalendarResponse, error) {
	//mock
	pcs := []*pb.PriceCalendar{}
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, -3)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(300),
	})
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, -2)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(200),
	})
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, -1)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(100),
	})
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, 3)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(1300),
	})
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, 2)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(1200),
	})
	pcs = append(pcs, &pb.PriceCalendar{
		FlightDatetime: timestamppb.New(req.FlightDatetime.AsTime().AddDate(0, 0, 1)),
		PassageType:    req.PassageType,
		Amount:         wrapperspb.Double(1100),
	})
	return &pb.PriceCalendarResponse{
		Result:        &common.Result{Code: "0"},
		PriceCalendar: pcs,
	}, nil
}
