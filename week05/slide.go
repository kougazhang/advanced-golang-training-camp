package main

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	PASS int = 1
	ERR  int = 2
)

//指标
type metrics struct {
	pass int
	err  int
}

type slideWindow struct {
	bucket int                //桶数量
	time   int64              //桶的创建时间
	event  map[int64]*metrics //保存事件
	data   *list.List         //保存窗口事件信息
	sync.RWMutex
}

// NewSlideWindow New 创建一个滑动窗口
func NewSlideWindow(bucket int) *slideWindow {
	sw := &slideWindow{}
	sw.bucket = bucket
	sw.data = list.New()
	return sw
}

// AddEvent metrics 中非零代表产生相应事件
func (sw *slideWindow) AddEvent(m metrics) {
	if m.pass!=0 {
		sw.incr(PASS)
	}
	if m.err!=0 {
		sw.incr(ERR)
	}


}

func (sw *slideWindow) incr(t int) {
	sw.Lock()
	defer sw.Unlock()
	nowTime := time.Now().Unix()
	if _, ok := sw.event[nowTime]; !ok {
		sw.event = make(map[int64]*metrics)
		sw.event[nowTime] = &metrics{}
	}else {

	}
	if sw.time == 0 {
		sw.time = nowTime
	}
	if sw.time != nowTime {
		sw.data.PushBack(sw.event[nowTime])
		sw.time = nowTime
		if sw.data.Len() > sw.bucket {
			//print("sw.data.Len()",sw.data.Len())
			for i := 0; i <= sw.data.Len()-sw.bucket; i++ {
				sw.data.Remove(sw.data.Front())
			}
		}
	}

	switch t {
	case PASS:
		sw.event[nowTime].pass++
	case ERR:
		sw.event[nowTime].err++
	default:
		log.Fatal("err type")
	}

}

// GetData 获取最新统计信息
func (sw *slideWindow)GetData() *metrics {
	sw.RLock()
	defer sw.RUnlock()
	var m = &metrics{}
	for i := sw.data.Front(); i != nil; i = i.Next() {
		v := i.Value.(*metrics)
		m.pass += v.pass
		m.err += v.err
	}
	for i := sw.data.Front(); i != nil; i = i.Next() {
		fmt.Print("data：",i.Value)
	}
	return m
}
