package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Conn represents a server connection.
type Conn struct {
	net.Conn
	Authenticated bool
}

// ListenAndServe starts a server at the specified address.
func ListenAndServe(
	host string, port int,
	handler func(conn *Conn, msg *Message, rd *AnyReaderWriter, w io.Writer) error,
	opened func(conn *Conn),
	closed func(conn *Conn),
	lnp *net.Listener,
) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	if lnp != nil {
		*lnp = ln
	}
	log.Printf("The server is now ready to accept connections on port %d\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		go handleConn(&Conn{Conn: conn}, handler, opened, closed)
	}
}

func handleConn(
	conn *Conn,
	handler func(conn *Conn, msg *Message, rd *AnyReaderWriter, w io.Writer) error,
	opened func(conn *Conn),
	closed func(conn *Conn),
) {

	opened(conn)
	defer closed(conn)

	defer conn.Close()

	rd := NewAnyReaderWriter(conn)

	for {
		msg, err := rd.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if msg != nil && msg.Command != "" {

			if msg.Command == "quit" {
				if msg.OutputType == RESP {
					io.WriteString(conn, "+OK\r\n")
				}
				return
			}
			err := handler(conn, msg, rd, conn)
			if err != nil {
				log.Println(err)
				return
			}

		} else {
			continue
		}

	}
}
