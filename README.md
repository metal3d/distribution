# Golang distributed calculation via RPC

That package aims to help you to create distributed calculation to nodes. It makes use of RPC and help you to auto-register nodes.

# Installation

You may install the package via

```
go get -u gopkg.in/metal3d/distribution.V1
```

# Usage

The given example will help you to understand how to create nodes, master and RPC endpoints.

At first, let's create a "tasks" package that implement a simple "Sum":

```golang
package tasks

import "time"

// a simple Arith type
type Arith int

// Implement a RPC endpoint. Keep in mind that methods should have 
// 2 arguments: one represents params, second represents reply and you should
// return an error
func (a *Arith) Sum(args *[]int, reply *int) error {
    time.Sleep(1 * time.Second) // simulate long process
    for _, v := range *args {
        *reply += v
    }
    return nil
}
```

Now let's create our "master" in "master/main.go":

```golang
package main

import (
    "fmt"
    "net/http"
    "net/rpc"
    dist "gopkg.in/metal3d/distribution.V1"
)

func main(){
    // handler a register endpoint for nodes
    dist.RegisterMaster()

    // create a test endpoint to call nodes
    http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request){
        args := []int{2,4,6}
        reply := 0
        err := dist.Call("Arith.Sum", &args, &reply)
        if err != nil {
            fmt.Log("Error: ", err)
        } else {
            fmt.Prinln("Reply from node", node.Addr, reply)
        }
    })
    fmt.Println("Listening :3000")
    http.ListenAndServe(":3000", nil)
}
```

That master node is now able to register nodes, the call to  `dist.RegisterMaster()` has registered a "register-node" handler. It listens on ":3000" port.

**NOTE** `RegisterMaster` will start a goroutine that will check if nodes are alive. When a node doesn't respond, so the master removes it from the list.

Now, create node source code in "node/main.go":

```golang
package main

import (
    "fmt"
    "net"
    "net/http"
    "net/rpc"
    "./tasks"
    dist "gopkg.in/metal3d/distribution.V1"
)

func main(){
    // now, register Arith as a RPC endpoint
    rpc.Register(new(tasks.Arith))

    // and let RPC to register our HTTP endpoints
    rpc.HandleHTTP()

    // open a tcp socket to serve node.
    // The ":0" will let os to give us a random port.
    l, err := net.Listen("tcp", ":0")
    if err != nil {
        log.Fatal(err)
    }
    
    // Register the node on master (localhost:3000).
    go dist.RegisterNode(l.Addr().String(), "localhost:3000")

    // use that socker throught http
    http.Serve(l, nil)
}
```

It's time to try:

```bash
# open a terminal
$ cd master
$ go run main.go
Listening :3000

# open another terminal
$ cd node
$ go run main.go

# open a third term
curl localhost:3000/sum
```

If you check terminal where you launched master, you'll see the reply from the node.

