package damianopetrungaro

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/damianopetrungaro/go-cache"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("damianopetrungaro/go-cache.empty", func(b *testing.B) {
		inmem := cache.NewInMemory[string, []byte](10*time.Second, 10_000)
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

	b.Run("damianopetrungaro/go-cache.prefilled", func(b *testing.B) {
		inmem := cache.NewInMemory[string, []byte](10*time.Second, 10_000)
		var k, v, ttl = "k", []byte("value"), time.Second

		for i := 0; i < 10_000; i++ {
			inmem.Set(context.Background(), fmt.Sprintf("%d", i), v, ttl)
			got, err := inmem.Get(context.Background(), fmt.Sprintf("%d", i))
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(got, v) {
				b.Fatal("not equal")
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(context.Background(), k)
				_ = inmem.Set(context.Background(), k, v, ttl)
				_ = inmem.Delete(context.Background(), k)
			}
		})
	})

	b.Run("damianopetrungaro/go-cache.prefilled_with_cleanup", func(b *testing.B) {
		inmem := cache.NewInMemory[string, []byte](10*time.Second, 10_000)
		var k, v, ttl = "k", []byte("value"), 500 * time.Millisecond

		for i := 0; i < 10_000; i++ {
			inmem.Set(context.Background(), fmt.Sprintf("%d", i), v, ttl)
		}

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
