package server

import (
	"bufio"
	"io"
	"strings"

	"github.com/junostorage/resp"
)

// Type is resp type
type Type int

const (
	Null Type = iota
	RESP
	Telnet
	HTTP
	JSON
)

// String return a string for type.
func (t Type) String() string {
	switch t {
	default:
		return "Unknown"
	case Null:
		return "Null"
	case RESP:
		return "RESP"
	case Telnet:
		return "Telnet"
	case HTTP:
		return "HTTP"
	case JSON:
		return "JSON"
	}
}

// Message is a resp message
type Message struct {
	Command    string
	Values     []resp.Value
	ConnType   Type
	OutputType Type
}

// AnyReaderWriter is resp or native reader writer.
type AnyReaderWriter struct {
	rd *bufio.Reader
	wr io.Writer
	ws bool
}

// NewAnyReaderWriter returns an AnyReaderWriter object.
func NewAnyReaderWriter(rd io.Reader) *AnyReaderWriter {
	ar := &AnyReaderWriter{}
	if rd2, ok := rd.(*bufio.Reader); ok {
		ar.rd = rd2
	} else {
		ar.rd = bufio.NewReader(rd)
	}
	if wr, ok := rd.(io.Writer); ok {
		ar.wr = wr
	}
	return ar
}

// ReadMessage reads the next resp message.
func (ar *AnyReaderWriter) ReadMessage() (*Message, error) {
	_, err := ar.rd.ReadByte()

	if err != nil {
		return nil, err
	}

	if err := ar.rd.UnreadByte(); err != nil {
		return nil, err
	}

	// MultiBulk also reads telnet
	return ar.readMultiBulkMessage()
}

func commandValues(values []resp.Value) string {
	if len(values) == 0 {
		return ""
	}
	return strings.ToLower(values[0].String())
}

func (ar *AnyReaderWriter) readMultiBulkMessage() (*Message, error) {
	rd := resp.NewReader(ar.rd)
	v, _, _, err := rd.ReadMultiBulk()

	if err != nil {
		return nil, err
	}
	values := v.Array()
	if len(values) == 0 {
		return nil, nil
	}

	return &Message{Command: commandValues(values), Values: values, ConnType: RESP, OutputType: RESP}, nil
}
