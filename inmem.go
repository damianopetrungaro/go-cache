package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var _ Cache[string, any] = &InMem[string, any]{}

type expiresAt int64

func (ea expiresAt) isExpired() bool {
	i := int64(ea)
	return time.Now().UnixNano() > i && i != int64(NoExpiration)
}

type item[V any] struct {
	val       V
	expiresAt expiresAt
}

// InMem is a Cache implementation which interacts with an in-memory map
// It is concurrent safe
type InMem[K comparable, V any] struct {
	items  map[K]item[V]
	mu     sync.Mutex
	ticker *time.Ticker
}

// NewInMemory returns a InMem instance
func NewInMemory[K comparable, V any](cleanUpInterval time.Duration) *InMem[K, V] {
	inmem := &InMem[K, V]{
		items:  map[K]item[V]{},
		ticker: time.NewTicker(cleanUpInterval),
	}

	go func() {
		for range inmem.ticker.C {
			inmem.cleanup()
		}
	}()

	return inmem
}

// Get retrieves an item from an in-memory map
func (i *InMem[K, V]) Get(ctx context.Context, key K) (V, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return *new(V), fmt.Errorf("%w: %s", ErrNotGet, ctx.Err())
	default:
	}
	item, ok := i.items[key]
	if !ok {
		return *new(V), ErrNotFound
	}

	if item.expiresAt.isExpired() {
		return *new(V), ErrExpired
	}

	return item.val, nil
}

// Set stores an item to an in-memory map
func (i *InMem[K, V]) Set(ctx context.Context, key K, val V, ttl time.Duration) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", ErrNotSet, ctx.Err())
	default:
	}
	exp := expiresAt(ttl)
	if ttl != NoExpiration {
		exp = expiresAt(time.Now().Add(ttl).UnixNano())
	}

	i.items[key] = item[V]{val: val, expiresAt: exp}
	return nil
}

// Delete removes an item to an in-memory map
func (i *InMem[K, V]) Delete(ctx context.Context, key K) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %s", ErrNotDelete, ctx.Err())
	default:
	}
	delete(i.items, key)
	return nil
}

// Close stops the inner ticker
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
