package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Fare struct {
	Id              int64
	OrgAirport      string
	ArrAirport      string
	PassageType     string
	FirstTravelDate time.Time
	LastTravelDate  time.Time
	Amount          float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type FareRepo interface {
	// db
	Pricing(ctx context.Context, fare *Fare) (*Fare, error)
	CreateFare(ctx context.Context, fare *Fare) (int64, error)
	UpdateFare(ctx context.Context, id int64, fare *Fare) error
	GetFare(ctx context.Context, id int64) (*Fare, error)
	DeleteFare(ctx context.Context, id int64) error
}

type FareUsecase struct {
	repo FareRepo
}

func NewFareUsecase(repo FareRepo, logger log.Logger) *FareUsecase {
	return &FareUsecase{repo: repo}
}

func (uc *FareUsecase) Pricing(ctx context.Context, fare *Fare) (p *Fare, err error) {
	p, err = uc.repo.Pricing(ctx, fare)
	return
}

func (uc *FareUsecase) Create(ctx context.Context, fare *Fare) (int64, error) {
	return uc.repo.CreateFare(ctx, fare)
}

func (uc *FareUsecase) Update(ctx context.Context, id int64, fare *Fare) error {
	return uc.repo.UpdateFare(ctx, id, fare)
}

func (uc *FareUsecase) Get(ctx context.Context, id int64) (*Fare, error) {
	return uc.repo.GetFare(ctx, id)
}

func (uc *FareUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteFare(ctx, id)
}
