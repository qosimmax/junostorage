package controller

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/junostorage/logger"

	"github.com/junostorage/controller/server"
	"github.com/junostorage/storage"

	"github.com/junostorage/resp"
)

var (
	errInvalidNumberOfArguments = errors.New("invalid number of arguments")
	logs                        *logrus.Logger
)

const (
	//@TODO feature read maxmem value from conf file
	maxmem = 1024 * 1024 * 1024 // 1GB
)

// Controller struct
type Controller struct {
	mu                     sync.RWMutex
	host                   string
	port                   int
	conns                  map[*server.Conn]bool
	statsTotalConns        int
	stopBackgroundExpiring bool
	stopWatchingMemory     bool
	outOfMemory            bool
	cache                  *storage.MemoryCache
}

func init() {
	logs = logger.GetLogger()
}

// ListenAndServe starts a new server
func ListenAndServe(host string, port int, httpPort int) error {
	return ListenAndServeEx(host, port, httpPort, nil)
}

// ListenAndServeEx function
func ListenAndServeEx(host string, port int, httpPort int, ln *net.Listener) error {

	c := &Controller{
		host:  host,
		port:  port,
		conns: make(map[*server.Conn]bool),
		cache: storage.New()}

	// watch memory
	go c.watchMemory()
	// expire checker
	go c.backgroundExpiring()

	defer func() {
		c.mu.Lock()
		c.stopBackgroundExpiring = true
		c.stopWatchingMemory = true
		c.mu.Unlock()
	}()

	handler := func(conn *server.Conn, msg *server.Message, rd *server.AnyReaderWriter, w io.Writer) error {

		err := c.handleInputCommand(conn, msg, w)
		if err != nil {
			logs.Error(err)
			return err
		}
		return nil
	}

	httpHandler := func(msg *server.Message, w http.ResponseWriter) error {

		err := c.handleInputCommand(nil, msg, w)
		if err != nil {
			logs.Error(err)
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

	//run http server
	go server.ListenHttpServer(host, httpPort, httpHandler)

	return server.ListenAndServe(host, port, handler, opened, closed, ln)
}

func (c *Controller) handleInputCommand(conn *server.Conn, msg *server.Message, w io.Writer) error {

	writeOutput := func(res string) error {
		switch msg.ConnType {
		default:
			err := fmt.Errorf("unsupported conn type: %v", msg.ConnType)
			return err

		case server.HTTP:
			fmt.Fprintf(w, res)

		case server.Telnet:
			_, err := io.WriteString(w, res)
			return err

		}
		return nil
	}
	// Ping. Just send back the response.
	if msg.Command == "ping" {
		switch msg.OutputType {
		case server.RESP:
			return writeOutput("+PONG\r\n")
		}
		return nil
	}

	writeErr := func(err error) error {
		switch msg.OutputType {
		case server.JSON:
			return writeOutput(fmt.Sprintf(`{"status":false, "error":"%v"}`, err))
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

	case storage.CmdSet,
		storage.CmdDel,
		storage.CmdHset,
		storage.CmdHdel,
		storage.CmdLpush,
		storage.CmdLpop,
		storage.CmdExpire:
		// write operations
		c.mu.Lock()
		defer c.mu.Unlock()

	case storage.CmdGet,
		storage.CmdKeys,
		storage.CmdHget,
		storage.CmdHgetAll,
		storage.CmdLindex,
		storage.CmdLlen:
		// read operations
		c.mu.RLock()
		defer c.mu.RUnlock()

	}

	res, err := c.command(msg, w)
	if err != nil {
		logs.Error(err)
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
	case storage.CmdGet:
		res, err = c.cmdGet(msg)

	case storage.CmdSet:
		res, err = c.cmdSet(msg)

	case storage.CmdKeys:
		res, err = c.cmdKeys(msg)

	case storage.CmdHset:
		res, err = c.cmdHset(msg)

	case storage.CmdHget:
		res, err = c.cmdHget(msg)

	case storage.CmdHgetAll:
		res, err = c.cmdHgetAll(msg)

	case storage.CmdHdel:
		res, err = c.cmdHdel(msg)

	case storage.CmdDel:
		res, err = c.cmdDel(msg)

	case storage.CmdLpush:
		res, err = c.cmdLpush(msg)

	case storage.CmdLindex:
		res, err = c.cmdLIndex(msg)

	case storage.CmdLlen:
		res, err = c.cmdLen(msg)

	case storage.CmdLpop:
		res, err = c.cmdLpop(msg)

	case storage.CmdExpire:
		res, err = c.cmdExpire(msg)

	}
	return
}

// backgroundExpiring watches for when items must expire from the cache
func (c *Controller) backgroundExpiring() {
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()

	for range t.C {
		c.mu.Lock()
		if c.stopBackgroundExpiring {
			c.mu.Unlock()
			return
		}

		keys := c.cache.ExpireList()
		for _, k := range keys {
			if c.cache.IsExpire(k) {
				c.cache.Del(k)
			}
		}
		c.mu.Unlock()

	}
}

func (c *Controller) watchMemory() {
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	var mem runtime.MemStats

	for range t.C {
		func() {
			c.mu.RLock()
			if c.stopWatchingMemory {
				c.mu.RUnlock()
				return
			}

			oom := c.outOfMemory
			c.mu.RUnlock()

			if oom {
				// runs a garbage collection
				runtime.GC()
			}
			runtime.ReadMemStats(&mem)

			c.mu.Lock()
			c.outOfMemory = int(mem.HeapAlloc) > maxmem
			c.mu.Unlock()

		}()
	}
}
