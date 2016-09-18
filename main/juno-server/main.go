package main

import (
	"flag"
	"log"

	"github.com/junostorage/controller"
)

var (
	port     int
	httpPort int
	host     string
)

func main() {

	flag.IntVar(&port, "p", 6380, "The listening port.")
	flag.IntVar(&httpPort, "http", 6382, "The http listening port.")
	flag.StringVar(&host, "h", "", "The listening host.")
	flag.Parse()

	if err := controller.ListenAndServe(host, port, httpPort); err != nil {
		log.Fatal(err)
	}
}
