package cache_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	. "github.com/damianopetrungaro/go-cache"
)

func TestInMem(t *testing.T) {

	t.Run("not found", func(t *testing.T) {
		inmem := newInMemHelper(t)

		val, err := inmem.Get(context.Background(), "one")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("find set value", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const k = "key"
		want := "value"
		if err := inmem.Set(context.Background(), k, want, NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		got, err := inmem.Get(context.Background(), k)
		if err != nil {
			t.Fatalf("could not get item: %s", err)
		}

		if got != want {
			t.Errorf("could not match value, got: %s. want:%s", got, want)
		}
	})

	t.Run("delete set value", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const k = "key"
		want := "value"
		if err := inmem.Set(context.Background(), k, want, NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Delete(context.Background(), k); err != nil {
			t.Fatalf("could not delete item: %s", err)
		}

		val, err := inmem.Get(context.Background(), k)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("concurrent set, get, and delete", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const c = 100
		wg := sync.WaitGroup{}
		wg.Add(c)

		for i := 0; i < 100; i++ {
			go func() {
				defer wg.Done()
				_, _ = inmem.Get(context.Background(), "one")
				_ = inmem.Set(context.Background(), "two", "two", time.Second)
				_ = inmem.Delete(context.Background(), "two")
			}()
		}

		wg.Wait()
	})

	t.Run("get expired value", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const k = "key"
		want := "value"
		if err := inmem.Set(context.Background(), k, want, time.Millisecond); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		time.Sleep(time.Millisecond)
		val, err := inmem.Get(context.Background(), k)
		if !errors.Is(err, ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("ensure cleanup behavior when expired item", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const longK, longV = "long key", "long value"
		const midK, midV = "mid key", "mid value"
		const shortK, shortV = "short key", "short value"
		const k, v = "key", "value"

		if err := inmem.Set(context.Background(), longK, longV, -5*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), midK, midV, -3*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), shortK, shortV, -1*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), k, v, 1*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if _, err := inmem.Get(context.Background(), longK); err == nil {
			t.Fatalf("could not find item: %s", longK)
		}

		if _, err := inmem.Get(context.Background(), midK); err == nil {
			t.Fatalf("could not find item: %s", midK)
		}

		if _, err := inmem.Get(context.Background(), shortK); err == nil {
			t.Fatalf("could find item: %s", shortK)
		}

		if found, _ := inmem.Get(context.Background(), k); found != v {
			t.Fatalf("could not find item: %s", k)
		}
	})

	t.Run("ensure cleanup behavior when no expired item", func(t *testing.T) {
		inmem := newInMemHelper(t)

		const longK, longV = "long key", "long value"
		const midK, midV = "mid key", "mid value"
		const shortK, shortV = "short key", "short value"
		const k, v = "key", "value"

		if err := inmem.Set(context.Background(), longK, longV, 5*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), midK, midV, 3*time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), shortK, shortV, time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Set(context.Background(), k, v, time.Minute); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if found, _ := inmem.Get(context.Background(), longK); found != longV {
			t.Fatalf("could not find item: %s", longK)
		}

		if found, _ := inmem.Get(context.Background(), midK); found != midV {
			t.Fatalf("could not find item: %s", midK)
		}

		if _, err := inmem.Get(context.Background(), shortK); err == nil {
			t.Fatalf("could find item: %s", shortK)
		}

		if found, _ := inmem.Get(context.Background(), k); found != v {
			t.Fatalf("could not find item: %s", k)
		}
	})
}

func newInMemHelper(t *testing.T) *InMem[string, string] {
	t.Helper()
	inmem := NewInMemory[string, string](time.Minute, 3)
	t.Cleanup(func() {
		if err := inmem.Close(); err != nil {
			t.Errorf("could not close inmem: %s", err)
		}
	})
	return inmem
}
