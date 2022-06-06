package redis_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/damianopetrungaro/go-cache/redis"
)

func Test_DefaultEncoder(t *testing.T) {
	data := []string{"One", "Two", "Three"}
	got, err := DefaultEncoder[any](data)
	if err != nil {
		t.Fatalf("could not encode item: %s", err)
	}

	want, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("could not encode item: %s", err)
	}

	if !bytes.Equal(got, want) {
		t.Error("could not match encoded item")
		t.Errorf("got: %v", got)
		t.Errorf("want: %v", want)
	}
}

func Test_DefaultDecoder(t *testing.T) {
	data := []byte(`["One", "Two", "Three"]`)

	var got []string
	if err := DefaultDecoder[any](data, &got); err != nil {
		t.Fatalf("could not encode item: %s", err)
	}

	var want []string
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatalf("could not encode item: %s", err)
	}

	if got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Error("could not match decoded item")
		t.Errorf("got: %v", got)
		t.Errorf("want: %v", want)
	}
}
