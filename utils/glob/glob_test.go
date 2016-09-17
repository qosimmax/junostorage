package glob

import "testing"

func TestMatch(t *testing.T) {

	values := []string{"hello", "hallo"}
	pattern := "h[ae]llo"

	for _, v := range values {
		matched, err := Match(pattern, v)
		if err != nil {
			t.Fatalf("match error:%v", err)
		}
		if !matched {
			t.Errorf("%v value not matched", v)
		}

	}

}
