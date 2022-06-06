package patrickn

import (
	"testing"

	"github.com/dgraph-io/ristretto"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("dgraph/ristretto", func(b *testing.B) {
		inmem, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7, // number of keys to track frequency of (10M).
			MaxCost:     1e7, // maximum cost of cache (1GB).
			BufferItems: 64,  // number of keys per Get buffer.
		})
		if err != nil {
			b.Fatal(err)
		}

		var k, v, cost = "k", []byte("value"), int64(1)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				inmem.Set(k, v, cost)
				inmem.Del(k)
			}
		})
	})
}
