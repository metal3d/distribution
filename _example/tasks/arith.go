package tasks

import (
	"net/rpc"
	"time"
)

// Arith overrides "int"
type Arith int

// Sum makes a sum of the whole intgers given as argument. It saves reponse in reply.
func (a *Arith) Sum(args *[]int, reply *int) error {
	time.Sleep(1 * time.Second) // to simulate long function
	*reply = 0
	for _, v := range *args {
		*reply += v
	}
	return nil
}

// RegisterArith is a simple function to register Arith structure as a RPC endpoint.
func RegisterArith() {
	rpc.Register(new(Arith))
}
