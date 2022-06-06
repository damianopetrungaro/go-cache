package patrickn

import (
	"fmt"
	"testing"
	"time"

	patrickmnCache "github.com/patrickmn/go-cache"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("patrickmn/go-cache.empty", func(b *testing.B) {
		inmem := patrickmnCache.New(5*time.Minute, 10*time.Minute)
		var k, v, ttl = "k", []byte("value"), time.Second
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				inmem.Set(k, v, ttl)
				inmem.Delete(k)
			}
		})
	})

	b.Run("patrickmn/go-cache.prefilled", func(b *testing.B) {
		inmem := patrickmnCache.New(5*time.Minute, 10*time.Minute)
		var k, v, ttl = "k", []byte("value"), time.Second

		for i := 0; i < 10_000; i++ {
			inmem.Set(fmt.Sprintf("%d", i), v, ttl)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				inmem.Set(k, v, ttl)
				inmem.Delete(k)
			}
		})
	})

	b.Run("patrickmn/go-cache.prefilled_with_cleanup", func(b *testing.B) {
		inmem := patrickmnCache.New(time.Second, time.Second)
		var k, v, ttl = "k", []byte("value"), 500 * time.Millisecond

		for i := 0; i < 10_000; i++ {
			inmem.Set(fmt.Sprintf("%d", i), v, ttl)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				inmem.Set(k, v, ttl)
				inmem.Delete(k)
			}
		})
	})
}
