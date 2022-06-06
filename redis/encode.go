package redis

import (
	"encoding/json"
)

// EncodeDecodeOption represents an Option which specify a strategy to encode and decode the items in the cache
func EncodeDecodeOption[K string, V any](enc Encoder[V], dec Decoder[*V]) Option[K, V] {
	return func(r *Redis[K, V]) {
		r.enc = enc
		r.dec = dec
		r.shouldEncodeDecode = true
	}
}

// Encoder represents a function used to encode an item as []byte to persist on redis
type Encoder[V any] func(val V) ([]byte, error)

// Decoder represents a function used to decode an item as []byte to persist on redis
type Decoder[V any] func(data []byte, val V) error

// DefaultEncoder is a default implementation of an Encoder. It transforms data to JSON.
func DefaultEncoder[V any](val V) ([]byte, error) {
	return json.Marshal(val)
}

// DefaultDecoder is a default implementation of a Decoder. It transforms data from JSON.
func DefaultDecoder[V any](data []byte, val V) error {
	if err := json.Unmarshal(data, val); err != nil {
		return err
	}

	return nil
}
