package distribution

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var masterDelay = time.Duration(5 * time.Second)
var masterDelayOnErr = time.Duration(500 * time.Millisecond)

// Contact the master continuously.
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
		cl := http.Client{}
		cl.Timeout = time.Duration(10 * time.Second)
		_, err := cl.Get(u.String())
		if err != nil {
			log.Println("Master", u.String(), "cannot be contacted, retrying in", masterDelayOnErr)
			time.Sleep(masterDelayOnErr)
			continue
		}
		time.Sleep(masterDelay)
	}

}

// RegisterNode starts the node listening process.
func RegisterNode(port string, master string) {
	// cleanup port
	port = strings.Replace(port, ":", "", -1)
	port = strings.Replace(port, "[", "", -1)
	port = strings.Replace(port, "]", "", -1)
	go registerNodeToMaster(port, master)
}
