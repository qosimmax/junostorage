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

func TestHTTPReader(t *testing.T) {

	data := "/Set/mykey/hallo/"
	buffer := bytes.NewBuffer([]byte(data))
	reader := NewAnyReaderWriter(buffer)

	_, err := reader.ReadHTTPMessage()
	if err != nil {
		t.Errorf("reader error:%v", err)
	}

}
