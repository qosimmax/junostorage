package storage

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	key := "My"
	value := "20"

	memcache := New()
	if err := memcache.Set(key, value); err != nil {
		t.Fatalf("Set error:%v", err)
	}

	got, err := memcache.Get(key)
	if err != nil {
		t.Fatalf("Get error:%v", err)
	}

	if value != got {
		t.Errorf("Want: %s, got: %s", value, got)
	}

}

func TestDictValues(t *testing.T) {
	key := "My"
	field := "age"
	value := "20"
	memcache := New()
	if err := memcache.HSet(key, field, value); err != nil {
		t.Fatalf("HSet error:%v", err)
	}

	if err := memcache.HSet(key, "name", "kasim"); err != nil {
		t.Fatalf("HSet error:%v", err)
	}

	got, err := memcache.HGet(key, field)
	if err != nil {
		t.Fatalf("HGet error:%v", err)
	}

	if value != got {
		t.Errorf("Want: %s, got: %s", value, got)
	}

	if _, err := memcache.HGetAll(key); err != nil {
		t.Fatalf("HGetAll error:%v", err)
	}

}

func TestListValues(t *testing.T) {
	key := "list"
	memcache := New()
	values := []string{"1", "2", "3"}

	if err := memcache.LPush(key, values...); err != nil {
		t.Fatalf("LPUSH error:%v", err)
	}

	if _, err := memcache.Llen(key); err != nil {
		t.Fatalf("LLEN error:%v", err)
	}

	if _, err := memcache.Lindex(key, -1); err != nil {
		t.Fatalf("LINDEX error:%v", err)
	}

	if _, err := memcache.LPop(key); err != nil {
		t.Fatalf("LPOP error:%v", err)
	}

}

func TestKeys(t *testing.T) {
	key := "storage"
	value := "redis"
	pattern := "*"

	memcache := New()
	if err := memcache.Set(key, value); err != nil {
		t.Fatalf("Set error:%v", err)
	}

	if _, err := memcache.Keys(pattern); err != nil {
		t.Fatalf("KEYS error:%v", err)
	}

}

func TestExpireKey(t *testing.T) {
	key := "storage"
	value := "redis"

	memcache := New()
	if err := memcache.Set(key, value); err != nil {
		t.Fatalf("Set error:%v", err)
	}

	if err := memcache.SetTTL(key, 60*time.Duration(time.Second)); err != nil {
		t.Fatalf("TTL error:%v", err)
	}

}
