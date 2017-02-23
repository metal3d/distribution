package distribution

import "net/rpc"

// Waiter is a struct that is returned by Go() method to be able to
// wait for a Node response. It handles the rpc.Call to be able to get
// errors if any.
type Waiter struct {
	Node    *Node
	rpcCall *rpc.Call
	client  *rpc.Client
}

// Wait for the response caller.
func (w *Waiter) Wait() {
	defer func(w *Waiter) {
		w.Node.Count--
		w.client.Close()
	}(w)
	<-w.rpcCall.Done
}

// Error returns the rpc.Call error if any.
func (w *Waiter) Error() error {
	return w.rpcCall.Error
}
