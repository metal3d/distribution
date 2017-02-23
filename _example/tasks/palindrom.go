package tasks

import (
	"fmt"
	"net/rpc"
)

type Range struct {
	Start int
	End   int
}

type Palindrom int

func IsPalindrom(value int) bool {
	bin := fmt.Sprintf("%b", value)
	l := len(bin)

	// get left part
	left := bin[:l/2]

	// givup the middle bit
	pad := 0
	if l%2 != 0 {
		pad = 1
	}
	// get the right part
	right := bin[l/2+pad:]

	// reverse right part
	n := len(right)
	runes := make([]rune, n)
	for _, r := range right {
		n--
		runes[n] = r
	}
	right = string(runes)

	return left == right
}

// Check if int "v" is binary palindrom.
func (p *Palindrom) Check(v *int, r *bool) error {
	*r = IsPalindrom(*v)
	return nil
}

// CheckN count the binary palindrom in a giver Range.
func (p *Palindrom) CheckN(v *Range, c *int) error {
	fmt.Println("Caluclate palindroms in range", v)
	for i := v.Start; i < v.End; i++ {
		if IsPalindrom(i) {
			*c++
		}
	}
	return nil
}

// RegisterPalindrom registers interface in RPC handlers.
func RegisterPalindrom(s *rpc.Server) {
	s.Register(new(Palindrom))
}
