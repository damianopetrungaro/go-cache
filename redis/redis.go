package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/damianopetrungaro/go-cache"
)

var _ cache.Cache[string, string] = &Redis[string, string]{}

// Option represent a function which applies changes to a Redis cache instance
type Option[K string, V any] func(*Redis[K, V])

// Redis is a cache.Cache implementation which interacts with a redis server
type Redis[K string, V any] struct {
	cl                 *redis.Client
	enc                Encoder[V]
	dec                Decoder[*V]
	shouldEncodeDecode bool
}

// New returns a Redis instance
func New[K string, V any](cl *redis.Client, opts ...Option[K, V]) *Redis[K, V] {
	r := &Redis[K, V]{
		cl: cl,
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

// Get retrieves an item from a redis server
func (r *Redis[K, V]) Get(ctx context.Context, k K) (V, error) {
	val := new(V)
	switch r.shouldEncodeDecode {
	case false:
		switch err := r.cl.Get(ctx, string(k)).Scan(val); {
		case err == redis.Nil:
			return *new(V), fmt.Errorf("%w:%s", cache.ErrNotFound, err)
		case err == nil:
			return *val, nil
		default:
			return *new(V), fmt.Errorf("%w:%s", cache.ErrNotGet, err)
		}
	default:
		data, err := r.cl.Get(ctx, string(k)).Bytes()
		switch {
		case err == redis.Nil:
			return *new(V), fmt.Errorf("%w:%s", cache.ErrNotFound, err)
		case err != nil:
			return *new(V), fmt.Errorf("%w:%s", cache.ErrNotGet, err)
		}

		if err := r.dec(data, val); err != nil {
			return *new(V), fmt.Errorf("%w:%s", cache.ErrNotGet, err)
		}
	}
	return *val, nil
}

// Set stores an item to a redis server
func (r *Redis[K, V]) Set(ctx context.Context, k K, v V, ttl time.Duration) error {

	switch r.shouldEncodeDecode {
	case false:
		if err := r.cl.Set(ctx, string(k), v, ttl).Err(); err != nil {
			return fmt.Errorf("%w:%s", cache.ErrNotSet, err)
		}
	default:
		data, err := r.enc(v)
		if err != nil {
			return fmt.Errorf("%w:%s", cache.ErrNotSet, err)
		}

		if err := r.cl.Set(ctx, string(k), data, ttl).Err(); err != nil {
			return fmt.Errorf("%w:%s", cache.ErrNotSet, err)
		}
	}

	return nil
}

// Delete removes an item from a redis server
func (r *Redis[K, V]) Delete(ctx context.Context, k K) error {
	if err := r.cl.Del(ctx, string(k)).Err(); err != nil {
		return fmt.Errorf("%w:%s", cache.ErrNotDelete, err)
	}
	return nil
}

// Close closes the redis connection
func (r *Redis[K, V]) Close() error {
	return r.cl.Close()
}
