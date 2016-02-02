package lc

import (
	"sync"
	"testing"
	"time"
)

func Test_SetAndGet(t *testing.T) {
	lc := NewLocalCopy(1*time.Millisecond, func(*LocalCopy) {})

	lc.Set("foo", "bar")

	if foo, found := lc.Get("foo"); !found {
		t.Errorf("Could not find ket (foo)")
	} else if f, ok := foo.(string); !ok {
		t.Errorf("Foo expected to be string")
	} else if f != "bar" {
		t.Errorf("Foo expected to be bar")
	}
}

func Test_Remove(t *testing.T) {
	lc := NewLocalCopy(1*time.Millisecond, func(*LocalCopy) {})

	lc.Set("foo", "bar")

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Could not find key (foo)")
	}

	lc.Remove("foo")

	if _, found := lc.Get("foo"); found {
		t.Errorf("Key (foo) should not exist anymore")
	}
}

func Test_Update(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	lc := NewLocalCopy(10*time.Millisecond, func(lc *LocalCopy) {
		lc.Set("foo", "bar")
		wg.Done()
	})

	if _, found := lc.Get("foo"); found {
		t.Errorf("Key (foo) should not exist yet")
	}

	wg.Wait()

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Key foo should exist!")
	}
}

func Test_UpdateImmediately(t *testing.T) {
	lc := NewImmediateLocalCopy(10*time.Millisecond, func(lc *LocalCopy) {
		lc.Set("foo", "bar")
	})

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Key foo should exist!")
	}
}
