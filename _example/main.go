package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/metal3d/distribution"
	"github.com/metal3d/distribution/_example/handlers"
	"github.com/metal3d/distribution/_example/tasks"
	//"gopkg.in/metal3d/distribution.v0/_example/handlers"
	//"gopkg.in/metal3d/distribution.v0/_example/tasts"
	//"gopkg.in/metal3d/distribution.v0"
)

var (
	masterAddr = ":3000"
	node       = false
	master     = fmt.Sprintf("%s:%d", "localhost", masterAddr)
	debug      = 1
)

func main() {
	flag.StringVar(&masterAddr, "addr", masterAddr, "addr to listen for master")
	flag.BoolVar(&node, "node", node, "declare this process as node")
	flag.StringVar(&master, "master", master, "master address for node")
	flag.IntVar(&debug, "debug", debug, "see logs")
	flag.Parse()

	distribution.Debug = debug

	if node {
		// register that node to the master
		node := distribution.RegisterNode(master)

		//register RPC endpoints
		tasks.RegisterArith(node.Server)
		tasks.RegisterPalindrom(node.Server)

		// generate endpoints
		node.HandleHTTP()

		log.Println("Starting node")
		log.Fatal(node.Serve())

	} else {

		// Serve master endpoint - note that the http server is
		// separated.
		go func() {
			log.Fatal(distribution.ServeMaster(masterAddr))
		}()

		// handlers to test RPC calls
		http.HandleFunc("/sum", handlers.Sum)
		http.HandleFunc("/palindrom", handlers.Palindrom)

		log.Println("Starting server on :3001")
		log.Fatal(http.ListenAndServe(":3001", nil))

	}
}
