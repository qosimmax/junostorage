package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func ListenHttpServer(host string, port int,
	httpHandler func(msg *Message, w http.ResponseWriter) error) error {

	bind := fmt.Sprintf("%v:%v", host, port)
	s := &http.Server{
		Addr:           bind,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", Handler(httpHandler))
	log.Printf("The http server listening port %d\n", port)
	return s.ListenAndServe()
}

func Handler(httpHandler func(msg *Message, w http.ResponseWriter) error) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("Content-Type", "application/json")

		buffer := bytes.NewBuffer([]byte(r.URL.Path))
		reader := NewAnyReaderWriter(buffer)
		msg, err := reader.ReadHTTPMessage()
		if err != nil {
			fmt.Fprintf(wr, `{"status":false, "error":"%v"}`, err)
			return
		}

		httpHandler(msg, wr)

	}
}
