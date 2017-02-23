package tasks

import (
	"fmt"
	"testing"
)

var checks = map[int]bool{
	9:  true,
	10: false,
	16: false,
	17: true,
}

func TestPalindrom(t *testing.T) {

	p := new(Palindrom)
	for v, res := range checks {
		fr := IsPalindrom(v)
		t.Log("Function result", v, fmt.Sprintf("%b", v), fr)
		if res != fr {
			t.Error(v, ":", fr, "should be", res)
		}
		var r bool
		val := v
		err := p.Check(&val, &r)
		if err != nil {
			t.Error(err)
		}
		t.Log("Palindrom type", val, fmt.Sprintf("%b", val), r)
		if r != res {
			t.Error(val, "Palindrom RPC function set to", r, "should be", res)
		}
	}

}
