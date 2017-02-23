package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/metal3d/distribution"
	//"gopkg.in/metal3d/distribution.v0"
)

// HandleSum will call Arith.Sum on one node.
func Sum(w http.ResponseWriter, r *http.Request) {
	args := []int{rand.Int(), rand.Int(), rand.Int()}
	reply := 0
	node, err := distribution.Call("Arith.Sum", &args, &reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%+v %s", node, err)))
		return
	}
	// Write response
	w.Write([]byte(fmt.Sprintf("Reponse from %s: %d", node.Addr, reply)))

}
