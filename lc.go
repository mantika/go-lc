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

func NewLocalCopy(interval time.Duration, updateFunc func(*LocalCopy)) *LocalCopy {
	lc := &LocalCopy{data: cmap.New(), interval: interval, updateFunc: updateFunc}
	lc.ticker = time.NewTicker(interval)
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
	return lc
}
