package allegro_test

import (
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
)

func BenchmarkLogger(b *testing.B) {
	b.Run("allegro/bigcache.empty", func(b *testing.B) {
		cfg := bigcache.DefaultConfig(10 * time.Minute)
		cfg.Verbose = false
		cfg.Logger = log.New(io.Discard, "", log.LstdFlags)
		inmem, err := bigcache.NewBigCache(cfg)
		if err != nil {
			b.Fatal(err)
		}

		var k, v = "k", []byte("value")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				_ = inmem.Set(k, v)
				_ = inmem.Delete(k)
			}
		})
	})

	b.Run("allegro/bigcache.prefilled", func(b *testing.B) {
		cfg := bigcache.DefaultConfig(10 * time.Minute)
		cfg.Verbose = false
		cfg.Logger = log.New(io.Discard, "", log.LstdFlags)
		inmem, err := bigcache.NewBigCache(cfg)
		if err != nil {
			b.Fatal(err)
		}

		var k, v = "k", []byte("value")

		for i := 0; i < 10_000; i++ {
			_ = inmem.Set(fmt.Sprintf("%d", i), v)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				_ = inmem.Set(k, v)
				_ = inmem.Delete(k)
			}
		})
	})

	b.Run("allegro/bigcache.prefilled_with_cleanup", func(b *testing.B) {
		cfg := bigcache.DefaultConfig(500 * time.Millisecond)
		cfg.Verbose = false
		cfg.Logger = log.New(io.Discard, "", log.LstdFlags)
		inmem, err := bigcache.NewBigCache(cfg)
		if err != nil {
			b.Fatal(err)
		}

		var k, v = "k", []byte("value")

		for i := 0; i < 10_000; i++ {
			_ = inmem.Set(fmt.Sprintf("%d", i), v)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = inmem.Get(k)
				_ = inmem.Set(k, v)
				_ = inmem.Delete(k)
			}
		})
	})
}
