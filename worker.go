package distribution

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
	"strings"
	"time"
)

type NodeServer struct {
	conn     *net.Conn
	listener net.Listener
	*rpc.Server
}

func (ns *NodeServer) HandleHTTP() {
	ns.Server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
}

func (ns *NodeServer) Serve() error {
	for {
		conn, err := ns.listener.Accept()
		if err != nil {
			return err
		}

		go ns.Server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
	return nil
}

// Contact the master continuously.
func registerNodeToMaster(master string) *NodeServer {

	server := rpc.NewServer()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	// register to master
	port := strings.Replace(l.Addr().String(), "[", "", -1)
	port = strings.Replace(port, "]", "", -1)
	port = strings.Replace(port, ":", "", -1)
	u := url.URL{}
	u.Host = master
	u.Scheme = "http"
	u.Path = RegisterEndpoint
	q := u.Query()
	q.Set("port", port)
	u.RawQuery = q.Encode()
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	fmt.Println("Contact master", u.String())
	_, err = client.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	node := NodeServer{
		Server:   server,
		listener: l,
	}
	return &node
}

// RegisterNode starts the node listening process.
func RegisterNode(master string) *NodeServer {
	// cleanup port
	return registerNodeToMaster(master)

}
