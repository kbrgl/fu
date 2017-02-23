package shallowradix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	tree := New()
	insertValues(t, tree)
	pfs := tree.Prefixes()
	assert.Equal(t, "github.com/", pfs[0])
	assert.Equal(t, "golang.org/pkg/path", pfs[1])
}

func insertValues(t *testing.T, tree *Tree) {
	shouldInsert(t, tree, "github.com/kbrgl/shallowradix")
	shouldInsert(t, tree, "github.com/kbrgl")
	shouldInsert(t, tree, "github.com/")
	shouldNotInsert(t, tree, "github.com/kbrgl/isnochys-syntax")
	shouldInsert(t, tree, "golang.org/pkg/path")
	shouldNotInsert(t, tree, "golang.org/pkg/path/filepath")
}

func shouldInsert(t *testing.T, tree *Tree, s string) {
	if !tree.Insert(s) {
		t.Errorf("val %s should have inserted into tree", s)
	}
}

func shouldNotInsert(t *testing.T, tree *Tree, s string) {
	if tree.Insert(s) {
		t.Errorf("val %s should not have inserted into tree", s)
	}
}
