package lc

import (
	"time"

	"github.com/streamrail/concurrent-map"
)

type LocalCopy struct {
	data       cmap.ConcurrentMap
	interval   time.Duration
	updateFunc func(*LocalCopy)
	quit       chan struct{}
	ticker     *time.Ticker
}

func (lc *LocalCopy) Get(key string) (interface{}, bool) {
	return lc.data.Get(key)
}

func (lc *LocalCopy) Set(key string, value interface{}) {
	lc.data.Set(key, value)
}

func (lc *LocalCopy) Remove(key string) {
	lc.data.Remove(key)
}

func (lc *LocalCopy) start() {
	lc.ticker = time.NewTicker(lc.interval)
	lc.quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-lc.ticker.C:
				lc.updateFunc(lc)
			case <-lc.quit:
				lc.ticker.Stop()
				return
			}
		}
	}()
}

func NewLocalCopy(interval time.Duration, updateFunc func(*LocalCopy)) *LocalCopy {
	lc := &LocalCopy{data: cmap.New(), interval: interval, updateFunc: updateFunc}
	lc.start()
	return lc
}
func NewImmediateLocalCopy(interval time.Duration, updateFunc func(*LocalCopy)) *LocalCopy {
	lc := &LocalCopy{data: cmap.New(), interval: interval, updateFunc: updateFunc}
	lc.updateFunc(lc)
	lc.start()
	return lc
}
