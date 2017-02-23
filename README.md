# Golang distributed calculation via RPC

That package aims to help you to create distributed calculation to nodes. It makes use of RPC and help you to auto-register nodes.

Distribution package ease the node registration and heartbeat. Master checks nodes, and nodes checks master. When a node is killed/stopped, so the master detects and removes that nodes. New nodes can register to master.

That package provides 2 methods to call RPC nodes:

- `distribution.Call(string, interface{}, interface{})` to make a sync call
- `distribution.Go(string, interface{}, interface{})` to make an async call

Both method do the same call, but the returned values are not used the same. First return an error, while second return a `*Waiter` that can be `nil` in case of error. The Waiter handles used `*Node` and a `Wait()` method bloking while the node has not answered.

Waiter handler Node and Wait() method.

So, `Go()` method is probably better if you want to know wich node answered. See the palindrom handler in `_example` directory to see a complex calcluation on several nodes.

Example:

```golang
node1, err := distribution.Call('Bayesian.GetClassification', &dataset1, &response1)
if err != nil {
    //error
}

node2, err := distribution.Call('Bayesian.GetClassification', &dataset2, &response2)
if err != nil {
    //error
}

// or

w1 = distribution.Go('Bayesian.GetClassification', &dataset1, &response1)
if w1 == nil {
    // error
}
w2 = distribution.Go('Bayesian.GetClassification', &dataset2, &response2)
if w2 == nil {
    //error
} 

// if no error, let's wait channels
<-w1.Wait()
<-w2.Wait()
// At this time, both rpc calls succeded, we can use responses
```

The `Go` call is probably what you will really need to make async calls. But keep in mind that you can also call `Call` method in goroutines.

# Installation

You may install the package via

```
go get -u gopkg.in/metal3d/distribution.v1
```

# Try the example

To not interfer with your packages and install example binary, please follow this instructions that export another GOPATH to temporary directory.

That example, that is in `_example` directory, will build a tiny docker image and launches master and node containers. You will be able to call `/sum` endpoint that will send a random sum on one node.

You can scale up and down node list, stop master, restart master, to see what happends.

Please, 

```bash

# install example without installing binary (-d)
go get -d -u gopkg.in/metal3d/distribution.v1/...
cd $GOPATH/src/gopkg.in/distribution.v1/_example
make build
docker-compose up

# open a new terminal and do:
# -> scale up nodes
docker-compose scale node=4

# try calculation
for i in $(seq 6); do curl -s localhost:10000/sum & done; wait

# -> scale node down
docker-compose scale node=2

# re-try calculation
for i in $(seq 6); do curl -s localhost:10000/sum & done; wait

# stop master to see that nodes will try to reconnect
docker-compose stop master

# after a while, try to restart master, nodes will
# be redetected.
docker-compose start master

# IMPORTANT
# then please stop docker containers and cleanup
docker-compose stop
docker-compose down -v

```


# Usage

The following example will help you to understand how to create nodes, master and RPC endpoints.

We will create a "master" that can handler "nodes" connections. The master will open port "3000".
We will implement a "/sum" endpoint that will call "node" and call `Arith.Sum`.

We will create a "node" listening on a random port. That node will register itself to "localhost:3000" that is the master. That registration will send the node listening port. So, you will be able to launch serveral nodes in different terminals. Begin with only one to be sure.

Let's create our "master" in "master/main.go":

```golang
package main

import (
    "fmt"
    "net/http"
    "net/rpc"
    dist "gopkg.in/metal3d/distribution.v1"
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

Now, create node source code in "node/main.go" with a `Arith` type that can respond to RPC calls:

```golang
package main

import (
    "fmt"
    "net"
    "time"
    "net/http"
    "net/rpc"
    dist "gopkg.in/metal3d/distribution.v1"
)

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

func main(){
    // now, register Arith as a RPC endpoint
    rpc.Register(new(Arith))

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

# afterward, you can open several terminals and launch other nodes.

# open a third term to call master/sum endpoint on port 3000:
curl localhost:3000/sum
```

If you check terminal where you launched master, you'll see the reply from the node.

