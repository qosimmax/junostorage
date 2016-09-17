package server

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {

	data := "SET mkey 10\r\n"
	buffer := bytes.NewBuffer([]byte(data))
	reader := NewAnyReaderWriter(buffer)

	_, err := reader.ReadMessage()
	if err != nil {
		t.Errorf("reader error:%v", err)
	}

}
