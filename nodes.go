package distribution

// Nodes is a sortable Node collection.
type nodes []*Node

// Len implements sort.Interface Len method.
func (n nodes) Len() int {
	return len(n)
}

// Less implements sort.Interface Less method. We sort nodes by their "Count" property
// that is the number of tasks being in progress.
func (n nodes) Less(i, j int) bool {
	return n[i].Count < n[j].Count
}

// Swap implements sort.Interface Swap method.
func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
