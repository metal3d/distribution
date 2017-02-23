package distribution

import (
	"errors"
	"fmt"
	"net/rpc"
	"sort"
	"time"
)

var (
	Nodes            = nodes{}
	RegisterEndpoint = "/register-node"
	Debug            = false
	HeartBeatDelay   = time.Duration(2 * time.Second)
)

// Node represent a calculation node.
type Node struct {
	Addr  string
	Port  string
	Count int
}

// Call makes a rpc call to the less used node. If no node can answer so an error is returned.
func Call(endpoint string, args interface{}, reply interface{}) (*Node, error) {
	sort.Sort(Nodes)
	for _, node := range Nodes {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", node.Addr, node.Port))
		if err != nil {
			continue
		}
		defer client.Close()
		node.Count++
		client.Call(endpoint, args, reply)
		node.Count--
		return node, nil
	}
	return nil, errors.New("No nodes can reply")
}

// Go makes an async call to the less used node. It returns rpc.Call with a Done property
// that is a chan to wait. If no nodes can be contacted, so the return value is nil.
func Go(endpoint string, args interface{}, reply interface{}) *Waiter {
	sort.Sort(Nodes)
	for _, node := range Nodes {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", node.Addr, node.Port))
		if err != nil {
			continue
		}
		node.Count++
		return &Waiter{
			Node:    node,
			client:  client,
			rpcCall: client.Go(endpoint, args, reply, nil),
		}
	}
	return nil
}
