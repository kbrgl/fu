// Package shallowradix implements a radix tree of depth 1.
package shallowradix

import "strings"

// Tree is like a radix tree except it only stores the first level of prefixes
// and no values.
type Tree struct {
	prefixes []string
}

// New returns a new Tree
func New() *Tree {
	return &Tree{}
}

// Insert inserts a string into the tree.
func (t *Tree) Insert(s string) bool {
	for i, prefix := range t.prefixes {
		if strings.HasPrefix(s, prefix) {
			return false
		} else if strings.HasPrefix(prefix, s) {
			t.prefixes[i] = s
			return true
		}
	}
	t.prefixes = append(t.prefixes, s)
	return true
}

// Prefixes returns all of the depth-1 nodes of the tree.
func (t Tree) Prefixes() []string {
	return t.prefixes
}
