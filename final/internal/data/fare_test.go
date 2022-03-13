package data

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/webmin7761/go-school/homework/final/internal/biz"
	"github.com/webmin7761/go-school/homework/final/internal/data/ent"
	"github.com/webmin7761/go-school/homework/final/internal/data/ent/enttest"
)

func TestEnt(t *testing.T) {
	// client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client := enttest.Open(t, "sqlite3", "./foo.db?cache=shared&mode=memory&_fk=1")
	defer client.Close()

	err := client.Schema.Create(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCreateFare(t *testing.T) {
	client, err := ent.Open("sqlite3", "./foo.db?cache=shared&mode=memory&_fk=1")
	defer client.Close()

	fareRepo := NewFareRepo(&Data{db: client}, log.NewStdLogger(os.Stdout))

	FirstTravelDate, _ := time.Parse("2006-01-02 15:04:05", "2018-04-23 12:24:51")

	// fare, err := fareRepo.Pricing(context.Background(), &biz.Fare{
	// 	OrgAirport:      "SHA",
	// 	ArrAirport:      "PEK",
	// 	PassageType:     "ADT",
	// 	FirstTravelDate: FirstTravelDate.Add(time.Hour * 24),
	// 	LastTravelDate:  FirstTravelDate.Add(time.Hour * 24),
	// })
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if fare != nil {
	// 	fmt.Println("有相同规则运价存在")
	// 	return
	// }
	// fmt.Printf("%#v\r", fare)

	for i := 0; i < 2; i++ {
		err = fareRepo.CreateFare(context.Background(), &biz.Fare{
			OrgAirport:      "SHA",
			ArrAirport:      "PEK",
			PassageType:     "ADT",
			FirstTravelDate: FirstTravelDate,
			LastTravelDate:  FirstTravelDate.Add(time.Hour * 24 * 30),
			Amount:          234.0,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
