package distribution

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// remove node when it doesn't respond. See hearBeatNodes().
func removeNode(node *Node) {
	for i, n := range Nodes {
		if n == node {
			if Debug > 0 {
				log.Println("Removing node", node.Addr)
			}
			Nodes = append(Nodes[:i], Nodes[i+1:]...)
		}
	}
}

// Start to serve. Server has a endpoint to register node and another to call "testSum".
func registerMasterHandler() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc(RegisterEndpoint, func(w http.ResponseWriter, req *http.Request) {
		a := strings.Split(req.RemoteAddr, ":")[0]
		remoteport := req.URL.Query().Get("port")
		for _, node := range Nodes {
			if node.Addr == a && remoteport == node.Port {
				if Debug > 1 {
					log.Println("Node already registered:", node)
				}
				w.WriteHeader(201)
				return
			}
		}
		Nodes = append(Nodes, &Node{a, remoteport, 0})
		if Debug > 0 {
			nlist := ""
			for _, n := range Nodes {
				nlist += fmt.Sprintf("%s:%s, ", n.Addr, n.Port)
			}
			log.Println("Nodes:", nlist)
		}
		w.WriteHeader(201)
	})
	return mux
}

// ServeMaster starts a http server on given address. That server handle a
// function to let nodes registering in stack.
// The endpoint for node registration is set in RegisterEndpoint var.
//
// TODO: set TLS capabilities
func ServeMaster(addr string) error {
	// Separate with a new server/muxer
	server := &http.Server{}
	server.Addr = addr
	server.Handler = registerMasterHandler()
	return server.ListenAndServe()
}
