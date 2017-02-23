package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/metal3d/distribution"
	"github.com/metal3d/distribution/_example/handlers"
	"github.com/metal3d/distribution/_example/tasks"
	//"gopkg.in/metal3d/distribution.v0/_example/handlers"
	//"gopkg.in/metal3d/distribution.v0/_example/tasts"
	//"gopkg.in/metal3d/distribution.v0"
)

var (
	port   = 10000
	node   = false
	master = fmt.Sprintf("%s:%d", "localhost", port)
	debug  = true
)

func main() {
	flag.IntVar(&port, "port", port, "port to listen")
	flag.BoolVar(&node, "node", node, "declare this process as node")
	flag.StringVar(&master, "master", master, "master address")
	flag.BoolVar(&debug, "debug", debug, "see logs")
	flag.Parse()

	distribution.Debug = debug

	if node {
		// open a listen interface
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}

		// register that node to the master
		go distribution.RegisterNode(l.Addr().String(), master)

		//register RPC endpoints
		tasks.RegisterArith()
		tasks.RegisterPalindrom()

		// generate endpoints
		rpc.HandleHTTP()

		// and serve !
		log.Fatal(http.Serve(l, nil))
	} else {
		distribution.RegisterMaster()

		// handlers to test RPC calls
		http.HandleFunc("/sum", handlers.Sum)
		http.HandleFunc("/palindrom", handlers.Palindrom)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
