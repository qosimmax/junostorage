package controller

import (
	"bytes"
	"testing"

	"github.com/junostorage/storage"

	"github.com/junostorage/controller/server"
)

var c = &Controller{cache: storage.New()}

func readMessage(data string) (message *server.Message, err error) {
	buffer := bytes.NewBuffer([]byte(data))
	reader := server.NewAnyReaderWriter(buffer)
	message, err = reader.ReadMessage()
	return
}

func TestCmdSet(t *testing.T) {

	data := "SET mkey 10\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdSet(message)
	if err != nil {
		t.Errorf("cmdSet error:%v", err)
	}

}

func TestCmdGet(t *testing.T) {

	data := "Get mkey\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdGet(message)
	if err != nil {
		t.Errorf("cmdGet error:%v", err)
	}

}

func TestCmdHset(t *testing.T) {

	data := "HSET person age 20\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdHset(message)
	if err != nil {
		t.Errorf("cmdHset error:%v", err)
	}

}

func TestCmdHGet(t *testing.T) {

	data := "HGET person age\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdHget(message)
	if err != nil {
		t.Errorf("cmdHGet error:%v", err)
	}

}

func TestCmdHGetAll(t *testing.T) {

	data := "HGETALL person\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdHgetAll(message)
	if err != nil {
		t.Errorf("cmdHGet error:%v", err)
	}

}

func TestCmdHDel(t *testing.T) {

	data := "HDEL person age\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdHdel(message)
	if err != nil {
		t.Errorf("cmdHDel error:%v", err)
	}

}

func TestCmdLpush(t *testing.T) {

	data := "LPUSH list 1 2 3 4 5\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdLpush(message)
	if err != nil {
		t.Errorf("cmdLpush error:%v", err)
	}

}

func TestCmdLindex(t *testing.T) {

	data := "LINDEX list 1\r\n"

	message, err := readMessage(data)
	if err != nil {
		t.Fatalf("reader error:%v", err)
	}

	_, err = c.cmdLIndex(message)
	if err != nil {
		t.Errorf("cmdLIndex error:%v", err)
	}

}
