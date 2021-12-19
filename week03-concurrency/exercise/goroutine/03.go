package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	tr := NewTracker()
	// tr.Run 是消费者, 消费 channel ch, 消费完毕后会向 channel stop 发一个信号
	go tr.Run()
	// channel ch 的生产者
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	// 对于 context 不太了解
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	// 关闭 channel ch
	tr.Shutdown(ctx)
}

func NewTracker() *Tracker {
	return &Tracker{
		ch:   make(chan string, 1),
		stop: make(chan struct{}),
	}
}

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func (t *Tracker) Event(ctx context.Context, data string) (err error) {
	select {
	case t.ch <- data:
		fmt.Println("[Event] send data ", data)
		return nil
	case <-ctx.Done():
		fmt.Println("[Event] done")
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(time.Second)
		fmt.Println("[Run] receive data ", data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	fmt.Println("shutdown ...")
	close(t.ch)
	// 如果 Tracker.Run 中已向 Tracker.stop 发信号
	select {
	// 那么会走 t.stop
	case <-t.stop:
		// 否则, 会走这里, 说明超时了
	case <-ctx.Done():
	}
}
