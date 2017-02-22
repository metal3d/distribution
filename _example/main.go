package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sort"

	"github.com/metal3d/distribution/_example/tasks"

	"github.com/metal3d/distribution"
)

var (
	port   = 10000
	node   = false
	master = fmt.Sprintf("%s:%d", "localhost", port)
)

func main() {
	flag.IntVar(&port, "port", port, "port to listen")
	flag.BoolVar(&node, "node", node, "declare this process as node")
	flag.StringVar(&master, "master", master, "master address")
	flag.Parse()

	distribution.Debug = true

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

		// generate endpoints
		rpc.HandleHTTP()

		// and serve !
		log.Fatal(http.Serve(l, nil))
	} else {
		distribution.RegisterMaster()

		// Add a simple test
		http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
			args := []int{rand.Int(), rand.Int(), rand.Int()}
			reply := 0
			// that sorts nodes from least used to higher used
			sort.Sort(distribution.Nodes)
			for _, node := range distribution.Nodes {
				client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", node.Addr, node.Port))
				if err != nil {
					continue // try another
				}
				defer client.Close()

				// Write response
				client.Call("Arith.Sum", &args, &reply)
				w.Write([]byte(fmt.Sprintf("Reponse: %d", reply)))

				// stop iteration
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error, no nodes has been contacted"))
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
