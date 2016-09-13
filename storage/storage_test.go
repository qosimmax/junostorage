package storage

import "testing"

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

func TestHSetHGet(t *testing.T) {
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

}

func TestLPush(t *testing.T) {
	key := "list"
	memcache := New()
	values := []string{"1", "2", "3"}
	if err := memcache.LPush(key, values...); err != nil {
		t.Fatalf("LPush error:%v", err)
	}

	if err := memcache.LPush(key, "4", "5", "6"); err != nil {
		t.Fatalf("LPush error:%v", err)
	}

}
