package client

import (
	"testing"
	"time"
)

func TestDialClinet(t *testing.T) {

	testCases := []struct {
		cmd  string
		args []interface{}
		res  string
	}{

		{
			cmd:  "set",
			args: []interface{}{"mykey", "hallo"},
			res:  `OK`,
		},

		{
			cmd:  "get",
			args: []interface{}{"mykey"},
			res:  `hallo`,
		},
		{
			cmd:  "hset",
			args: []interface{}{"person", "name", "nemo"},
			res:  "1",
		},
		{
			cmd:  "hset",
			args: []interface{}{"person", "age", "20"},
			res:  "1",
		},
		{
			cmd:  "expire",
			args: []interface{}{"mykey", 10},
			res:  "1",
		},
		{
			cmd:  "del",
			args: []interface{}{"mylist"},
			res:  "0",
		},
		{
			cmd:  "lpush",
			args: []interface{}{"mylist", "hi", 1, "test", 4.5},
			res:  "4",
		},
		{
			cmd:  "llen",
			args: []interface{}{"mylist"},
			res:  "4",
		},
		{
			cmd:  "lindex",
			args: []interface{}{"mylist", 0},
			res:  "4.5",
		},
		{
			cmd:  "lpop",
			args: []interface{}{"mylist"},
			res:  "4.5",
		},
		{
			cmd:  "del",
			args: []interface{}{"mylist"},
			res:  "1",
		},
		{
			cmd:  "del",
			args: []interface{}{"mykey"},
			res:  "1",
		},
	}

	con, err := DialTimeout("localhost:6380", 60*time.Duration(time.Second))
	if err != nil {
		t.Fatalf("Dial error:%v", err)
	}

	for _, testCase := range testCases {
		val, err := con.Do(testCase.cmd, testCase.args...)
		if err != nil {
			t.Errorf("On cmd send error:%v", err)
		}

		if val.Error() != nil {
			t.Errorf("resp error:%v", val.Error())
		}

		if testCase.res != val.String() {
			t.Errorf("Expected the result to be `%v`, but instead found it to be `%v`, cmd:%v",
				testCase.res, val, testCase.cmd)
		}

	}

}
