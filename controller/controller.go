package controller

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/junostorage/controller/server"
	"github.com/junostorage/storage"

	"github.com/junostorage/resp"
)

var (
	errInvalidNumberOfArguments = errors.New("invalid number of arguments")
)

// Controller struct
type Controller struct {
	mu              sync.RWMutex
	host            string
	port            int
	conns           map[*server.Conn]bool
	statsTotalConns int
	cache           *storage.MemoryCache
}

// ListenAndServe starts a new server
func ListenAndServe(host string, port int) error {
	return ListenAndServeEx(host, port, nil)
}

// ListenAndServeEx function
func ListenAndServeEx(host string, port int, ln *net.Listener) error {

	c := &Controller{
		host:  host,
		port:  port,
		conns: make(map[*server.Conn]bool),
		cache: storage.New()}

	//run expire checker
	go c.backgroundExpiring()

	handler := func(conn *server.Conn, msg *server.Message, rd *server.AnyReaderWriter, w io.Writer) error {

		err := c.handleInputCommand(conn, msg, w)
		if err != nil {
			return err
		}
		return nil
	}

	opened := func(conn *server.Conn) {
		c.mu.Lock()
		c.conns[conn] = true
		c.statsTotalConns++
		c.mu.Unlock()
	}
	closed := func(conn *server.Conn) {
		c.mu.Lock()
		delete(c.conns, conn)
		c.mu.Unlock()
	}
	return server.ListenAndServe(host, port, handler, opened, closed, ln)
}

func (c *Controller) handleInputCommand(conn *server.Conn, msg *server.Message, w io.Writer) error {

	writeOutput := func(res string) error {
		switch msg.ConnType {
		default:
			err := fmt.Errorf("unsupported conn type: %v", msg.ConnType)
			return err

		case server.RESP:
			_, err := io.WriteString(w, res)
			return err

		}
	}
	// Ping. Just send back the response. No need to put through the pipeline.
	if msg.Command == "ping" {
		switch msg.OutputType {
		case server.RESP:
			return writeOutput("+PONG\r\n")
		}
		return nil
	}

	writeErr := func(err error) error {
		switch msg.OutputType {
		case server.RESP:
			if err == errInvalidNumberOfArguments {
				return writeOutput("-ERR wrong number of arguments for '" + msg.Command + "' command\r\n")
			}
			v, _ := resp.ErrorValue(errors.New("ERR " + err.Error())).MarshalRESP()
			return writeOutput(string(v))
		}

		return nil
	}

	// choose the locking strategy
	switch msg.Command {
	default:
		c.mu.RLock()
		defer c.mu.RUnlock()
	case "set", "del", "hset", "hgetall", "hdel", "lpop", "lpush", "expire":
		// write operations
		c.mu.Lock()
		defer c.mu.Unlock()

	case "get", "hget", "lindex", "llen":
		// read operations
		c.mu.RLock()
		defer c.mu.RUnlock()

	}

	res, err := c.command(msg, w)
	if err != nil {
		return writeErr(err)
	}

	if res != "" {
		if err := writeOutput(res); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) command(msg *server.Message, w io.Writer) (res string, err error) {
	switch msg.Command {
	default:
		err = fmt.Errorf("unknown command '%s'", msg.Values[0])
	case "get":
		res, err = c.cmdGet(msg)

	case "set":
		res, err = c.cmdSet(msg)

	case "hset":
		res, err = c.cmdHset(msg)

	case "hget":
		res, err = c.cmdHget(msg)

	case "hgetall":
		res, err = c.cmdHgetAll(msg)

	case "hdel":
		res, err = c.cmdHdel(msg)

	case "del":
		res, err = c.cmdDel(msg)

	case "lpush":
		res, err = c.cmdLpush(msg)

	case "lindex":
		res, err = c.cmdLIndex(msg)

	case "llen":
		res, err = c.cmdLen(msg)

	case "lpop":
		res, err = c.cmdLpop(msg)

	case "expire":
		res, err = c.cmdExpire(msg)

	}
	return
}

// backgroundExpiring watches for when items must expire from the cache
func (c *Controller) backgroundExpiring() {
	for {
		c.mu.Lock()
		keys := c.cache.ExpireList()

		for _, k := range keys {
			if c.cache.IsExpire(k) {
				c.cache.Del(k)
			}
		}
		c.mu.Unlock()
		time.Sleep(time.Second / 3)
	}
}
