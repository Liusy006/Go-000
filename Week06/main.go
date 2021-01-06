package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Bucket struct {
	second int

	success int32
	failed  int32
	timeout int32
	refuse  int32
}

type Counter struct {
	slots []Bucket

	index int32
	size  int32

	stopChan chan struct{}
}

func NewCounter(num int32) *Counter {
	return &Counter{
		slots:    make([]Bucket, num),
		index:    0,
		size:     num,
		stopChan: make(chan struct{}),
	}
}

func (c *Counter) Start() {
	c.slots[c.index].second = time.Now().Second()

	ticker := time.NewTicker(time.Second)
	isRunning := true
	for {
		select {
		case <-ticker.C:
			for {
				old := c.index
				new := (old + 1) % c.size
				swap := atomic.CompareAndSwapInt32(&c.index, old, new)
				if swap {
					//清空当前槽的下一个槽的数据
					next := (old + 1) % c.size
					c.slots[next].success = 0
					c.slots[next].failed = 0
					c.slots[next].timeout = 0
					c.slots[next].refuse = 0
					c.slots[next].second = time.Now().Second() + 1
					break
				}
			}

		case <-c.stopChan:
			isRunning = false
		}

		if !isRunning {
			break
		}
	}
}

func (c *Counter) Stop() {
	c.stopChan <- struct{}{}
}

func (c *Counter) AddSuccess() int32 {
	return atomic.AddInt32(&c.slots[c.index].success, 1)
}

func (c *Counter) AddFailed() int32 {
	return atomic.AddInt32(&c.slots[c.index].failed, 1)
}

func (c *Counter) AddTimeout() int32 {
	return atomic.AddInt32(&c.slots[c.index].timeout, 1)
}

func (c *Counter) AddRefuse() int32 {
	return atomic.AddInt32(&c.slots[c.index].refuse, 1)
}

func (c *Counter) GetData() []Bucket {
	return c.slots
}

func main() {
	c := NewCounter(10)
	go c.Start()

	i := 0
	for {
		time.Sleep(time.Millisecond)
		c.AddFailed()

		if i > 10000 {
			break
		}
		i++
	}

	fmt.Print(c.GetData())
	c.Stop()
}
