package damianopetrungaro

import (
	"context"
	"testing"
	"time"

	inmemoryCache "github.com/damianopetrungaro/go-cache/inmem"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("damianopetrungaro/go-cache", func(b *testing.B) {
		inmem := inmemoryCache.New[string, []byte](5 * time.Minute)
		var k, v, ttl = "k", []byte("value"), time.Second
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(context.Background(), k)
				_ = inmem.Set(context.Background(), k, v, ttl)
				_ = inmem.Delete(context.Background(), k)
			}
		})
	})
}
