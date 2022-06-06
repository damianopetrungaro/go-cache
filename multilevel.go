package cache

import (
	"context"
	"time"
)

// MultiLevel is a Cache implementation which allow a multi level usage cache
type MultiLevel[K comparable, V any] struct {
	caches []Cache[K, V]
}

// NewMultiLevel returns a MultiLevel
func NewMultiLevel[K comparable, V any](cs ...Cache[K, V]) *MultiLevel[K, V] {
	return &MultiLevel[K, V]{
		caches: cs,
	}
}

// Get traverse all the caches, if all of them fail it returns a generic ErrNotGet
func (m *MultiLevel[K, V]) Get(ctx context.Context, k K) (V, error) {
	for _, c := range m.caches {
		val, err := c.Get(ctx, k)
		if err == nil {
			return val, nil
		}
	}

	return *new(V), ErrNotGet
}

// Set traverse all the caches, if all of them fail it returns a generic ErrNotSet
func (m *MultiLevel[K, V]) Set(ctx context.Context, k K, v V, ttl time.Duration) error {
	var succeed int

	for _, c := range m.caches {
		if err := c.Set(ctx, k, v, ttl); err == nil {
			succeed++
		}
	}

	if succeed == 0 {
		return ErrNotSet
	}

	return nil
}

// Delete traverse all the caches, if all of them fail it returns a generic ErrNotDelete
func (m *MultiLevel[K, V]) Delete(ctx context.Context, k K) error {
	var succeed int

	for _, c := range m.caches {
		if err := c.Delete(ctx, k); err == nil {
			succeed++
		}
	}

	if succeed == 0 {
		return ErrNotDelete
	}

	return nil
}

// Close closes all the caches, it returns the last failing error
func (m *MultiLevel[K, V]) Close() error {
	var closeErr error

	for _, c := range m.caches {
		if err := c.Close(); err != nil {
			closeErr = err
		}
	}

	return closeErr
}
