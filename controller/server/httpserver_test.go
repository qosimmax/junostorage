package server

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPServer(t *testing.T) {

	testCases := []struct {
		url string
		res string
	}{

		{
			url: "http://localhost:6382/set/mkey/hallo",
			res: `{"status":true}`,
		},

		{
			url: "http://localhost:6382/get/mkey",
			res: `{"status":true, "value":"hallo"}`,
		},

		{
			url: "http://localhost:6382/hset/person/name/nemo",
			res: `{"status":true}`,
		},

		{
			url: "http://localhost:6382/hset/person/age/25",
			res: `{"status":true}`,
		},

		{
			url: "http://localhost:6382/hgetall/person",
			res: `{"status":true, "value":[name nemo age 25]}`,
		},

		{
			url: "http://localhost:6382/hdel/person/age",
			res: `{"status":true, "value":1}`,
		},

		{
			url: "http://localhost:6382/expire/mkey/10",
			res: `{"status":true, "value":1}`,
		},

		{
			url: "http://localhost:6382/lpush/list/1/2/3",
			res: `{"status":true, "value":3}`,
		},

		{
			url: "http://localhost:6382/llen/list",
			res: `{"status":true, "value":3}`,
		},

		{
			url: "http://localhost:6382/lindex/list/1",
			res: `{"status":true, "value":2}`,
		},

		{
			url: "http://localhost:6382/lpop/list",
			res: `{"status":true, "value":3}`,
		},

		{
			url: "http://localhost:6382/del/mkey",
			res: `{"status":true, "value":"1"}`,
		},
	}

	// Iterating over the test cases
	for _, testCase := range testCases {

		resp, err := http.Get(testCase.url)
		if err != nil {
			t.Errorf("http error:%v, url:%v", err, testCase.url)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		result := string(body)
		if testCase.res != string(body) {
			t.Errorf("Expected the result to be `%v`, but instead found it to be `%v`,url:%v",
				testCase.res, result, testCase.url)
		}
	}

}
