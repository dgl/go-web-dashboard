package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dgl/go-web-dashboard/clients"
	"github.com/dgl/go-web-dashboard/ui"
)

var flagListenAddr = flag.String("listen", ":4000", "Listen address ([address]:port)")

func main() {
	flag.Parse()

	ui.New(clients.New())

	log.Printf("Starting HTTP server on %v", *flagListenAddr)
	log.Fatal(http.ListenAndServe(*flagListenAddr, nil))
}
