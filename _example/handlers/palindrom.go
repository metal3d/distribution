package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/metal3d/distribution"
	"github.com/metal3d/distribution/_example/tasks"
	//"gopkg.in/metal3d/distribution.v0/_example/tasks"
	//"gopkg.in/metal3d/distribution.v0"
)

// Palindrom will call several nodes to check palindrom from 0 to
// a number giver in "n" parameters. Eg. curl -s master/palindrom?n=10000
//
// That handler make use of Go() method, and records results in a "counters" slice.
// After nodes responses, the total number of binary palindrom is written in the
// caller response.
func Palindrom(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("n")
	log.Println("Checking number of palindroms for", v)
	ival, err := strconv.Atoi(v)
	if err != nil {
		log.Println("Error", err)
		return
	}

	// cut palindrom calculation to several part of 100 value
	// to be tested
	// counters will handler the entire results
	counters := []*int{}
	waiters := []*distribution.Waiter{}
	j := 0
	for i := 0; i <= ival; i += 101 {
		end := i + 100
		if end >= ival {
			end = ival
		}

		// append a result pointer in the counters slice
		counters = append(counters, new(int))

		// call Palindrom.CheckN with
		c := distribution.Go("Palindrom.CheckN", &tasks.Range{i, end}, counters[j])
		waiters = append(waiters, c)
		j++
	}
	// Now, wait for nodes
	for i := 0; i < len(waiters); i++ {
		waiter := waiters[i]
		fmt.Println("Waiting", waiter.Node.Addr)
		waiter.Wait()
		fmt.Println("Reponse", waiter.Node.Addr)
	}

	// sum results that are kept in counters
	sum := 0
	for _, v := range counters {
		sum += *v
	}

	// write the response
	resp := []byte(fmt.Sprintf("Palindrom counter to %d gives %d\n", ival, sum))
	w.Write(resp)
}
