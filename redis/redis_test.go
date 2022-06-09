package redis_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/damianopetrungaro/go-cache"
	. "github.com/damianopetrungaro/go-cache/redis"
)

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test")
	}

	options, err := redis.ParseURL(getRedisUriHelper(t))
	if err != nil {
		t.Fatal(err)
	}

	testHelper(
		t,
		New[string, string](redis.NewClient(options)),
	)

	testHelper(
		t,
		New[string, string](
			redis.NewClient(options),
			EncodeDecodeOption[string, string](DefaultEncoder[string], DefaultDecoder[*string]),
		),
	)
}

func testHelper(t *testing.T, redisCache *Redis[string, string]) {
	t.Helper()
	t.Run("not found", func(t *testing.T) {
		val, err := redisCache.Get(context.Background(), uuid.New().String())
		if !errors.Is(err, cache.ErrNotFound) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("find set value", func(t *testing.T) {
		var k = uuid.New().String()
		want := "value"
		if err := redisCache.Set(context.Background(), k, want, cache.NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		got, err := redisCache.Get(context.Background(), k)
		if err != nil {
			t.Fatalf("could not get item: %s", err)
		}

		if got != want {
			t.Errorf("could not match value, got: %s. want:%s", got, want)
		}
	})

	t.Run("delete set value", func(t *testing.T) {
		var k = uuid.New().String()
		want := "value"
		if err := redisCache.Set(context.Background(), k, want, cache.NoExpiration); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		if err := redisCache.Delete(context.Background(), k); err != nil {
			t.Fatalf("could not delete item: %s", err)
		}

		val, err := redisCache.Get(context.Background(), k)
		if !errors.Is(err, cache.ErrNotFound) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})

	t.Run("concurrent set, get, and delete", func(t *testing.T) {
		const c = 100
		wg := sync.WaitGroup{}
		wg.Add(c)

		for i := 0; i < 100; i++ {
			go func() {
				defer wg.Done()
				_, _ = redisCache.Get(context.Background(), uuid.New().String())
				_ = redisCache.Set(context.Background(), uuid.New().String(), "two", time.Second)
				_ = redisCache.Delete(context.Background(), uuid.New().String())
			}()
		}

		wg.Wait()
	})

	t.Run("get expired value", func(t *testing.T) {
		var k = uuid.New().String()
		want := "value"
		if err := redisCache.Set(context.Background(), k, want, time.Millisecond); err != nil {
			t.Fatalf("could not set item: %s", err)
		}

		time.Sleep(100 * time.Millisecond)
		val, err := redisCache.Get(context.Background(), k)
		if !errors.Is(err, cache.ErrNotGet) {
			t.Errorf("could not match not found error. got: %s", err)
		}

		if val != "" {
			t.Errorf("could not match default value, got: %s", val)
		}
	})
}

func getRedisUriHelper(t *testing.T) string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           log.New(io.Discard, "", log.LstdFlags),
	})
	if err != nil {
		t.Fatalf("could not start redis container: %s", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("could not get redis mapped port: %s", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("could not get redis host: %s", err)
	}

	return fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
}
