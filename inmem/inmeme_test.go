package inmem_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/damianopetrungaro/go-cache"
	. "github.com/damianopetrungaro/go-cache/inmem"
)

func TestInMem(t *testing.T) {

	t.Run("not found", func(t *testing.T) {
		inmem := newInMemHelper(t)

		val, err := inmem.Get(context.Background(), "one")
		if !errors.Is(err, cache.ErrNotFound) {
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
		if err := inmem.Set(context.Background(), k, want, cache.NoExpiration); err != nil {
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
		if err := inmem.Set(context.Background(), k, want, cache.NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := inmem.Delete(context.Background(), k); err != nil {
			t.Fatalf("could not delete item: %s", err)
		}

		val, err := inmem.Get(context.Background(), k)
		if !errors.Is(err, cache.ErrNotFound) {
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
		if !errors.Is(err, cache.ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})
}

func newInMemHelper(t *testing.T) *InMem[string, string] {
	inmem := New[string, string](time.Millisecond)
	t.Cleanup(func() {
		if err := inmem.Close(); err != nil {
			t.Errorf("could not close inmem: %s", err)
		}
	})
	return inmem
}
