package distribution

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var masterDelay = time.Duration(500 * time.Millisecond)

// Start to serve node. That function can be used to register RPC
// endpoints. See "tasks.RegisterArith()" as example that is
// called in that function.
// That function also call master endpoint to register current node giving
// used address and port.
func registerNodeToMaster(port, master string) {
	// record
	u := url.URL{}
	u.Scheme = "http"
	u.Host = master
	u.Path = RegisterEndpoint
	query := u.Query()
	query.Set("port", port)
	u.RawQuery = query.Encode()

	// call master to register this node
	for {
		// continuously try to connect to master
		if _, err := http.Get(u.String()); err == nil {
			return
		}
		log.Println("Master", u.String(), "cannot be contacted, retrying in", masterDelay)
		time.Sleep(masterDelay)
	}

}

// RegisterNode starts the node listening process.
func RegisterNode(port string, master string) {
	// cleanup port
	port = strings.Replace(port, ":", "", -1)
	port = strings.Replace(port, "[", "", -1)
	port = strings.Replace(port, "]", "", -1)
	registerNodeToMaster(port, master)
}
