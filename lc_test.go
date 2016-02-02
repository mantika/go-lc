package lc

import (
	"testing"
	"time"
)

func Test_SetAndGet(t *testing.T) {
	lc := NewLocalCopy(1*time.Millisecond, func(*Handler) {})

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
	lc := NewLocalCopy(1*time.Millisecond, func(*Handler) {})

	lc.Set("foo", "bar")

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Could not find key (foo)")
	}

	lc.Remove("foo")

	if _, found := lc.Get("foo"); found {
		t.Errorf("Key (foo) should not exist anymore")
	}
}

func Test_Fill(t *testing.T) {
	lc := NewLocalCopy(10*time.Millisecond, func(h *Handler) {
		h.Clean()
		h.Set("foo", "bar")
	})

	if _, found := lc.Get("foo"); found {
		t.Errorf("Key (foo) should not exist yet")
	}
	lc.Set("foo1", "bar")

	time.Sleep(20 * time.Millisecond)

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Key foo should exist!")
	}
	if _, found := lc.Get("foo1"); found {
		t.Errorf("Key (foo1) should not exist!")
	}
}

func Test_FillImmediately(t *testing.T) {
	lc := NewImmediateLocalCopy(10*time.Millisecond, func(h *Handler) {
		h.Set("foo", "bar")
	})

	if _, found := lc.Get("foo"); !found {
		t.Errorf("Key foo should exist!")
	}
}
