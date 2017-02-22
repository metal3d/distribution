package distribution

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// remove node when it doesn't respond. See hearBeatNodes().
func removeNode(node *Node) {
	for i, n := range Nodes {
		if n == node {
			Nodes = append(Nodes[:i], Nodes[i+1:]...)
		}
	}
}

// Checking nodes.
func hearBeatNodes() {
	for {
		if HeartBeatDelay < time.Duration(2*time.Second) {
			// we need to have (at this time) a highier delay to let
			// timeout to be one second shorter.
			//
			// Task: find another way to get a timeout lesser that hreatbeat delay.
			HeartBeatDelay = time.Duration(2 * time.Second)
		}
		time.Sleep(HeartBeatDelay)
		for _, node := range Nodes {
			go checkNode(node, HeartBeatDelay-time.Duration(1*time.Second))
		}
	}
}

// Check on node. If connection breaks or is timeout so remove this node.
func checkNode(node *Node, timeout time.Duration) {
	if Debug {
		log.Println("Checking", node.Addr)
	}
	// try to contact node
	c, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%s", node.Addr, node.Port),
		timeout)

	if err != nil {
		if Debug {
			log.Println("Removing node", node)
		}
		removeNode(node)
		return
	}
	c.Close()
}

// Start to serve. Server has a endpoint to register node and another to call "testSum".
func registerMasterHandler() {
	http.HandleFunc(RegisterEndpoint, func(w http.ResponseWriter, req *http.Request) {
		a := strings.Split(req.RemoteAddr, ":")[0]
		remoteport := req.URL.Query().Get("port")
		Nodes = append(Nodes, &Node{a, remoteport, 0})
		if Debug {
			fmt.Println("Nodes:", Nodes)
		}
		w.WriteHeader(201)
	})
}

// RegisterMaster register the master node endpoint that is able to
// get node registration. It also starts heartbeating nodes.
func RegisterMaster() {
	go hearBeatNodes()
	registerMasterHandler()
}
