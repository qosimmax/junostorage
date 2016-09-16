package client

import (
	"testing"
	"time"
)

func TestDialClinet(t *testing.T) {
	con, err := DialTimeout("localhost:6380", 60*time.Duration(time.Second))
	if err != nil {
		t.Fatalf("Dial error:%v", err)
	}
	val, err := con.Do("SET", "mykey", "111")
	if err != nil {
		t.Fatalf("On cmd send error:%v", err)
	}

	if val.Error() != nil {
		t.Fatalf("resp error:%v", err)
	}

}
