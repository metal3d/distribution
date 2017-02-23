# Golang distributed calculation via RPC

That package aims to help you to create distributed calculation to nodes. It makes use of RPC and help you to auto-register nodes.

Distribution package ease the node registration and node down detection. Master checks nodes registration, and removes it if it fails to contact. When a node is killed/stopped, so the master detects and removes that nodes. New nodes can register to master at any time.

To be able to make use of this package, you have to:

- create a master using `ServeMaster(interface)` where "interface" is a string like ":3000". You probably will call that method in a goroutine.
- create nodes using `node := RegisterNode(masterUrl)` where "masterUrl" is a string like "master.url:3000", "localhost:3000" or via IP "172.17.1.1:3000"
- make node service running via `node.Serve()`

Note: `node.Serve()` and `RegisterNode()` are blocking, if you need to continue process after that calls, you will need to handle them in a goroutine. 

That package provides 2 methods to call RPC nodes:

- `distribution.Call(string, interface{}, interface{})` to make a **sync** calls
- `distribution.Go(string, interface{}, interface{})` to make an **async** calls

Both methods does the same call despite the fact that `Go()` will not block process. 
Also,returned values are not used the same way !. 

`Call()` will return a `*Node` and an error, while `Go()` will return a `*Waiter` that can be `nil` in case of error. 

The Waiter handles contacted `*Node` and a `Wait()` method bloking while the node has not answered. 

Because `Go()` method is async, there is no way (at this time) to be sure that the routine is ok unless checking `Waiter.Error()`. 

Example:

```golang

// Sync calls
node1, err := distribution.Call('Bayesian.GetClassification', &dataset1, &response1)
if err != nil {
    //error
}

node2, err := distribution.Call('Bayesian.GetClassification', &dataset2, &response2)
if err != nil {
    //error
}

// Async calls

w1 = distribution.Go('Bayesian.GetClassification', &dataset1, &response1)
if w1.Error() != nil {
    // error
}
w2 = distribution.Go('Bayesian.GetClassification', &dataset2, &response2)
if w2.Error() != nil {
    //error
} 

// if no error, let's wait channels
<-w1.Wait()
<-w2.Wait()
// At this time, both rpc calls succeded, we can use responses
```

The `Go()` call is probably what you will really need to make async calls. But keep in mind that **you can also call `Call` method in goroutines**.

# Installation

You may install the package via

```
go get -u gopkg.in/metal3d/distribution.v1
```

# Try the example

That example, that is in `_example` directory, will build a tiny docker image and launches master and node containers. You will be able to call `/sum` and `/palindrom` endpoints.

`/sum` is only a routine that makes the sum of 3 random integers. 

`/palindrom?n=X` will calculate how many binary palindrom exists from 0 to "X". That example split calculation in several ranges that are sent to several nodes. It reduce result after the all nodes reply.

You can scale up and down node list and see that stopped nodes are detected and removed from the stack.

```bash

# install example without installing binary (-d)
$ go get -d -u gopkg.in/metal3d/distribution.v1/...
$ cd $GOPATH/src/gopkg.in/distribution.v1/_example
$ make

# open a new terminal and do:
$ cd $GOPATH/src/gopkg.in/distribution.v1/_example
# scale up nodes
$ docker-compose scale node=4

# try calculation
$ curl -s localhost:3001/sum
Response from 172.17.1.2: 345365445465464
$ curl -s localhost:3001/palindrom?n=2000
Palindrom counter to 2000: 92

# -> scale down nodes
$ docker-compose scale node=2

# re-try calculation
$ curl -s localhost:3001/sum
Response from 172.17.1.2: 734466677764353
$ curl -s localhost:3001/palindrom?n=2000
Palindrom counter to 2000: 92

# IMPORTANT
# then please stop docker containers and cleanup
$ make clean
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
    dist "gopkg.in/metal3d/distribution.v1"
)

func main(){

    // at first, let's register this process as
    // master. So that node can contact master
    // to register using master-ip:3000
    go func() {
        dist.ServeMaster(":3000")
    }()

    // This is optionnal, but for testing purpose we will
    // create a test endpoint to call nodes. So we create a HTTP server
    // listening on :3001
    http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request){
        args := []int{2, 4, 6}
        reply := 0
        err := dist.Call("Arith.Sum", &args, &reply)
        if err != nil {
            fmt.Log("Error: ", err)
        } else {
            fmt.Prinln("Reply from node", node.Addr, reply)
            w.Write(fmt.Sprintf("Response: %d\n", reply))
        }
    })
    fmt.Println("Listening :3001")
    http.ListenAndServe(":3001", nil)
}
```

That master node is now able to register nodes, the call to  `dist.ServeMaster()` has registered a "register-node" handler. It listens on ":3000" port. Nodes will be configured to hit that endpoint.

Now, create node source code in "node/main.go" with a `Arith` type that can respond to RPC calls:

```golang
package main

import (
    "fmt"
    "time"
    "net/http"
    dist "gopkg.in/metal3d/distribution.v1"
)

// a simple Arith type to handle RPC methods.
type Arith int


// sum values.
func sum(values []int) int {
    tot := 0
    for _, v := range valuse {
        tot += v
    }
    return tot
}

// Implement a RPC endpoint. Keep in mind that methods should have 
// 2 arguments: one represents params, second represents reply and you should
// return an error
func (a *Arith) Sum(args *[]int, reply *int) error {
    time.Sleep(1 * time.Second) // simulate long process
    *reply = a.sum(*args)
    return nil
}

func main(){

    // create a node server - port is the RPC master.
    // Note that "localhost:3000" is the RPC Master url.
    // If master is not on localhost, replace that url.
    node := dist.RegisterNode("localhost:3000")

    // now, register Arith as a RPC endpoint
    node.Server.Register(new(Arith))

    // and now handler requests
    node.HandleHTTP()

    // And start to serve that node.
    log.Println("Node listening")
    node.Serve()
}
```

It's time to try:

```bash
# open a terminal
$ cd master
$ go run main.go
Listening :3001

# open another terminal
$ cd node
$ go run main.go

# afterward, you can open several terminals and launch other nodes.

# open a third term to call master/sum endpoint on port 3000:
curl localhost:3001/sum
Response: 12
```

If you check terminal where you launched master, you'll see the reply from the node.

Keep in mind that `:3000` was declared to serve RPC call, and `:3001` to serve `sum` endpoint. It's not mandatory to use http to serve endpoints to call RPC. That's an example to show you how a master can send requests to nodes and get back results.

Check `_example` directory to see a more complex example.

