package main

import (
	"flag"
	"log"

	"github.com/junostorage/controller"
)

var (
	port int
	host string
)

func main() {

	flag.IntVar(&port, "p", 6380, "The listening port.")
	flag.StringVar(&host, "h", "", "The listening host.")
	flag.Parse()

	if err := controller.ListenAndServe(host, port); err != nil {
		log.Fatal(err)
	}
}
