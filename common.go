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

// Nodes is a sortable Node collection.
type nodes []*Node

// Len implements sort.Interface Len method.
func (n nodes) Len() int {
	return len(n)
}

// Less implements sort.Interface Less method. We sort nodes by their "Count" property
// that is the number of tasks being in progress.
func (n nodes) Less(i, j int) bool {
	return n[i].Count < n[j].Count
}

// Swap implements sort.Interface Swap method.
func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Call makes a rpc call to the less used node. If no node can answer so an error is returned.
func Call(endpoint string, args interface{}, reply interface{}) error {
	sort.Sort(Nodes)
	for _, node := range Nodes {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", node.Addr, node.Port))
		if err != nil {
			continue
		}
		client.Call(endpoint, args, reply)
		return nil
	}
	return errors.New("No nodes can reply")
}

// Go makes an async call to the less used node. It returns rpc.Call with a Done property
// that is a chan to wait. If no nodes can be contacted, so the return value is nil.
func Go(endpoint string, args interface{}, reply interface{}) *rpc.Call {
	sort.Sort(Nodes)
	for _, node := range Nodes {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", node.Addr, node.Port))
		if err != nil {
			continue
		}
		return client.Go(endpoint, args, reply, nil)
	}
	return nil
}
