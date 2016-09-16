package controller

import (
	"time"

	"github.com/junostorage/controller/server"
	"github.com/junostorage/storage"

	"github.com/junostorage/resp"
)

func (c *Controller) cmdGet(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 2 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	value, err := c.cache.Get(key)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.NullValue().MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.StringValue(value)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}

		return string(data), nil
	}

	return "", nil
}

func (c *Controller) cmdSet(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	value := msg.Values[2].String()

	if err = c.cache.Set(key, value); err != nil {
		return
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.SimpleStringValue("OK")
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil

	}

	return
}

func (c *Controller) cmdDel(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 2 {
		err = errInvalidNumberOfArguments
		return
	}

	val := 0
	key := msg.Values[1].String()
	if ok := c.cache.Del(key); ok {
		val = 1
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.IntegerValue(val)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil

	}

	return "", nil
}

func (c *Controller) cmdHset(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 4 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	field := msg.Values[2].String()
	value := msg.Values[3].String()

	if err = c.cache.HSet(key, field, value); err != nil {
		return
	}
	switch msg.OutputType {
	case server.RESP:
		oval := resp.IntegerValue(1)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil

	}

	return
}

func (c *Controller) cmdHget(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	field := msg.Values[2].String()

	value, err := c.cache.HGet(key, field)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.NullValue().MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	oval := resp.StringValue(value)

	switch msg.OutputType {
	case server.RESP:
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}

func (c *Controller) cmdHgetAll(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 2 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()

	values, err := c.cache.HGetAll(key)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.NullValue().MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	vals := make([]resp.Value, 0, 2)
	for _, v := range values {
		vals = append(vals, resp.StringValue(v))
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.ArrayValue(vals)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}

func (c *Controller) cmdHdel(msg *server.Message) (res string, err error) {

	if len(msg.Values) < 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	fields := []string{}
	for _, v := range msg.Values[2:] {
		fields = append(fields, v.String())
	}

	n, err := c.cache.HDel(key, fields...)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.IntegerValue(0).MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.IntegerValue(n)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}

func (c *Controller) cmdLpush(msg *server.Message) (res string, err error) {

	if len(msg.Values) < 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	list := []string{}
	for _, v := range msg.Values[2:] {
		list = append(list, v.String())
	}

	err = c.cache.LPush(key, list...)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.IntegerValue(0).MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		n, _ := c.cache.Llen(key)
		oval := resp.IntegerValue(n)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil

	}

	return
}

func (c *Controller) cmdLen(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 2 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()

	n, err := c.cache.Llen(key)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.IntegerValue(0).MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.IntegerValue(n)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil

	}

	return
}

func (c *Controller) cmdLIndex(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	index := msg.Values[2].Integer()

	value, err := c.cache.Lindex(key, index)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.NullValue().MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.StringValue(value)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}

func (c *Controller) cmdLpop(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 2 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()

	value, err := c.cache.LPop(key)
	if err != nil {

		if err == storage.ErrNullValue {
			data, _ := resp.NullValue().MarshalRESP()
			return string(data), nil
		}

		return "", err
	}

	switch msg.OutputType {
	case server.RESP:
		oval := resp.StringValue(value)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}

func (c *Controller) cmdExpire(msg *server.Message) (res string, err error) {

	if len(msg.Values) != 3 {
		err = errInvalidNumberOfArguments
		return
	}

	key := msg.Values[1].String()
	value := msg.Values[2].Integer()

	c.cache.SetTTL(key, time.Duration(value)*time.Second)

	switch msg.OutputType {
	case server.RESP:
		oval := resp.IntegerValue(1)
		data, err := oval.MarshalRESP()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return
}
