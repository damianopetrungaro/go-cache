package patrickn

import (
	"testing"
	"time"

	patrickmnCache "github.com/patrickmn/go-cache"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("patrickmn/go-cache", func(b *testing.B) {
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
}
