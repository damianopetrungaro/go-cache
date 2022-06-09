package cache

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestMultiLevel(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		multiLvl := newMultiLevel(t)

		val, err := multiLvl.Get(context.Background(), "one")
		if !errors.Is(err, ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("find set value", func(t *testing.T) {
		multiLvl := newMultiLevel(t)

		const k = "key"
		want := "value"
		if err := multiLvl.Set(context.Background(), k, want, NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		got, err := multiLvl.Get(context.Background(), k)
		if err != nil {
			t.Fatalf("could not get item: %s", err)
		}

		if got != want {
			t.Errorf("could not match value, got: %s. want:%s", got, want)
		}
	})

	t.Run("delete set value", func(t *testing.T) {
		multiLvl := newMultiLevel(t)

		const k = "key"
		want := "value"
		if err := multiLvl.Set(context.Background(), k, want, NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := multiLvl.Delete(context.Background(), k); err != nil {
			t.Fatalf("could not delete item: %s", err)
		}

		val, err := multiLvl.Get(context.Background(), k)
		if !errors.Is(err, ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("concurrent set, get, and delete", func(t *testing.T) {
		multiLvl := newMultiLevel(t)

		const c = 100
		wg := sync.WaitGroup{}
		wg.Add(c)

		for i := 0; i < 100; i++ {
			go func() {
				defer wg.Done()
				_, _ = multiLvl.Get(context.Background(), "one")
				_ = multiLvl.Set(context.Background(), "two", "two", time.Second)
				_ = multiLvl.Delete(context.Background(), "two")
			}()
		}

		wg.Wait()
	})

	t.Run("get expired value", func(t *testing.T) {
		multiLvl := newMultiLevel(t)

		const k = "key"
		want := "value"
		if err := multiLvl.Set(context.Background(), k, want, time.Millisecond); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		time.Sleep(time.Millisecond)
		val, err := multiLvl.Get(context.Background(), k)
		if !errors.Is(err, ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})
}

func newMultiLevel(t *testing.T) *MultiLevel[string, string] {
	t.Helper()
	inmem1 := NewInMemory[string, string](time.Second, 5)
	inmem2 := NewInMemory[string, string](time.Second, 5)
	multiLlv := NewMultiLevel[string, string](inmem1, inmem2)
	t.Cleanup(func() {
		if err := inmem1.Close(); err != nil {
			t.Errorf("could not close inmem1 level: %s", err)
		}
		if err := inmem2.Close(); err != nil {
			t.Errorf("could not close inmem2 level: %s", err)
		}
	})
	return multiLlv
}
