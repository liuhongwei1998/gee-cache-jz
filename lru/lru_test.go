package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("jz1", String("jz1"))
	if v, ok := lru.Get("jz1"); !ok || string(v.(String)) != "jz1" {
		t.Fatal("cache hit jz1 = jz1 failed")
	}
	if _, ok := lru.Get("jz2"); ok {
		t.Fatal("cache miss jz2 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)

	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
