package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	second int64

	success int32
	failed  int32
	timeout int32
	refuse  int32
}

func (b *Bucket) clear() {
	b.second = time.Now().Unix()
	b.success = 0
	b.failed = 0
	b.timeout = 0
	b.refuse = 0
}

func (b *Bucket) AddSuccess() int32 {
	return atomic.AddInt32(&b.success, 1)
}

func (b *Bucket) AddFailed() int32 {
	return atomic.AddInt32(&b.failed, 1)
}

func (b *Bucket) AddTimeout() int32 {
	return atomic.AddInt32(&b.timeout, 1)
}

func (b *Bucket) AddRefuse() int32 {
	return atomic.AddInt32(&b.refuse, 1)
}

type Counter struct {
	slots []*Bucket
	mu    sync.Mutex

	index int32
	size  int32

	stopChan chan struct{}
}

func NewCounter(num int32) *Counter {
	counter := &Counter{
		slots:    make([]*Bucket, 0, num),
		index:    0,
		size:     num,
		stopChan: make(chan struct{}),
	}

	for i := 0; i < int(num); i++ {
		counter.slots = append(counter.slots, &Bucket{})
	}

	return counter
}

func (c *Counter) Start() {
	c.slots[0].clear()
	ticker := time.NewTicker(time.Second)
	isRunning := true
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			c.index = (c.index + 1) % c.size
			c.slots[c.index].clear()
			c.mu.Unlock()

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

func (c *Counter) GetCurrentBucket() *Bucket {
	c.mu.Lock()
	defer c.mu.Unlock()
	b := c.slots[c.index]
	return b
}

func (c *Counter) GetData() []Bucket {
	d := make([]Bucket, 0, c.size)
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, value := range c.slots {
		d = append(d, *value)
	}
	return d
}

func main() {
	c := NewCounter(10)
	go c.Start()

	i := 0
	for {
		time.Sleep(time.Millisecond)
		b := c.GetCurrentBucket()
		b.AddFailed()

		if i > 10000 {
			break
		}
		i++
	}

	fmt.Print(c.GetData())
	c.Stop()
}
