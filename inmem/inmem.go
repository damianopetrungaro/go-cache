package inmem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damianopetrungaro/go-cache"
)

var _ cache.Cache[string, any] = &InMem[string, any]{}

type expiresAt int64

func (ea expiresAt) isExpired() bool {
	if time.Now().UnixNano() > int64(ea) && ea != -1 {
		return true
	}

	return false
}

type item[V any] struct {
	val       V
	expiresAt expiresAt
}

type InMem[K comparable, V any] struct {
	items  map[K]item[V]
	mu     sync.Mutex
	ticker *time.Ticker
}

func New[K comparable, V any](xxx time.Duration) *InMem[K, V] {
	inmem := &InMem[K, V]{
		items:  map[K]item[V]{},
		ticker: time.NewTicker(xxx),
	}

	go func() {
		for range inmem.ticker.C {
			inmem.cleanup()
		}
	}()

	return inmem
}

func (i *InMem[K, V]) Get(ctx context.Context, key K) (V, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return *new(V), fmt.Errorf("%w: %s", cache.ErrNotGet, ctx.Err())
	default:
	}
	item, ok := i.items[key]
	if !ok {
		return *new(V), cache.ErrNotFound
	}

	if item.expiresAt.isExpired() {
		return *new(V), cache.ErrExpired
	}

	return item.val, nil
}

func (i *InMem[K, V]) Set(ctx context.Context, key K, val V, ttl time.Duration) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", cache.ErrNotSet, ctx.Err())
	default:
	}
	exp := expiresAt(ttl)
	if ttl != cache.NoExpiration {
		exp = expiresAt(time.Now().Add(ttl).UnixNano())
	}

	i.items[key] = item[V]{val: val, expiresAt: exp}
	return nil
}

func (i *InMem[K, V]) Delete(ctx context.Context, key K) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", cache.ErrNotDelete, ctx.Err())
	default:
	}
	delete(i.items, key)
	return nil
}

func (i *InMem[K, V]) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.ticker.Stop()
	return nil
}

func (i *InMem[K, V]) cleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()
	for k, item := range i.items {
		if item.expiresAt.isExpired() {
			delete(i.items, k)
		}
	}
}
