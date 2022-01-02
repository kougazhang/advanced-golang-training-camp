package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Metrics struct {
}

func main() {

	sw := NewSlideWindow(2)
	rand.Seed(time.Now().UnixNano())
	var r int
	go func() {

		for i := 0; i < 1000; i++ {
			r = rand.Intn(3)
			if r == 1 {
				sw.AddEvent(metrics{pass: PASS, err: 0})
			} else {
				sw.AddEvent(metrics{pass: 0, err: ERR})
			}

			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		}
	}()

	for  {
		time.Sleep(time.Second)
		fmt.Println("统计：",sw.GetData())
	}
}
