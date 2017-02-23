package distribution

import (
	"errors"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"sort"
	"time"
)

var (
	Nodes            = nodes{}
	RegisterEndpoint = "/register-node"
	Debug            = 0
	HeartBeatDelay   = time.Duration(1 * time.Second)
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
		conn, err := net.DialTimeout("tcp",
			fmt.Sprintf("%s:%s", node.Addr, node.Port),
			HeartBeatDelay)
		if err != nil {
			removeNode(node)
			continue
		}
		client := jsonrpc.NewClient(conn)
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
		conn, err := net.DialTimeout("tcp",
			fmt.Sprintf("%s:%s", node.Addr, node.Port),
			HeartBeatDelay)
		if err != nil {
			removeNode(node)
			continue
		}
		client := jsonrpc.NewClient(conn)
		node.Count++
		return &Waiter{
			Node:    node,
			client:  client,
			rpcCall: client.Go(endpoint, args, reply, nil),
		}
	}
	return nil
}
