package lc

import (
	"sync"
	"time"
)

type LocalCopy struct {
	data     map[string]interface{}
	interval time.Duration
	fillFunc func(*Handler)
	ticker   *time.Ticker
	lock     sync.RWMutex
}

type Handler struct {
	lc *LocalCopy
}

func (h *Handler) Get(key string) (interface{}, bool) {
	i, o := h.lc.data[key]
	return i, o
}
func (h *Handler) Set(key string, value interface{}) {
	h.lc.data[key] = value
}
func (h *Handler) Clean() {
	h.lc.data = map[string]interface{}{}
}

func (lc *LocalCopy) Get(key string) (interface{}, bool) {
	lc.lock.RLock()
	i, o := lc.data[key]
	lc.lock.RUnlock()
	return i, o
}

func (lc *LocalCopy) Set(key string, value interface{}) {
	lc.lock.Lock()
	lc.data[key] = value
	lc.lock.Unlock()
}

func (lc *LocalCopy) Remove(key string) {
	lc.lock.Lock()
	delete(lc.data, key)
	lc.lock.Unlock()
}

func (lc *LocalCopy) start() {
	lc.ticker = time.NewTicker(lc.interval)
	go func() {
		for {
			select {
			case <-lc.ticker.C:
				lc.fill()
			}
		}
	}()
}

func (lc *LocalCopy) fill() {
	lc.lock.Lock()
	lc.fillFunc(&Handler{lc: lc})
	lc.lock.Unlock()
}

func NewLocalCopy(interval time.Duration, fillFunc func(*Handler)) *LocalCopy {
	lc := &LocalCopy{data: map[string]interface{}{}, interval: interval, fillFunc: fillFunc}
	lc.start()
	return lc
}
func NewImmediateLocalCopy(interval time.Duration, fillFunc func(*Handler)) *LocalCopy {
	lc := &LocalCopy{data: map[string]interface{}{}, interval: interval, fillFunc: fillFunc}
	lc.fill()
	lc.start()
	return lc
}
