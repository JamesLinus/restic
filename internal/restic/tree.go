package restic

import (
	"fmt"
	"sort"

	"github.com/restic/restic/internal/errors"

	"github.com/restic/restic/internal/debug"
)

// Tree is an ordered list of nodes.
type Tree struct {
	Nodes []*Node `json:"nodes"`
}

// NewTree creates a new tree object.
func NewTree() *Tree {
	return &Tree{
		Nodes: []*Node{},
	}
}

func (t Tree) String() string {
	return fmt.Sprintf("Tree<%d nodes>", len(t.Nodes))
}

// Equals returns true if t and other have exactly the same nodes.
func (t Tree) Equals(other *Tree) bool {
	if len(t.Nodes) != len(other.Nodes) {
		debug.Log("tree.Equals(): trees have different number of nodes")
		return false
	}

	for i := 0; i < len(t.Nodes); i++ {
		if !t.Nodes[i].Equals(*other.Nodes[i]) {
			debug.Log("tree.Equals(): node %d is different:", i)
			debug.Log("  %#v", t.Nodes[i])
			debug.Log("  %#v", other.Nodes[i])
			return false
		}
	}

	return true
}

// Insert adds a new node at the correct place in the tree.
func (t *Tree) Insert(node *Node) error {
	pos, _, err := t.binarySearch(node.Name)
	if err == nil {
		return errors.New("node already present")
	}

	// https://code.google.com/p/go-wiki/wiki/SliceTricks
	t.Nodes = append(t.Nodes, &Node{})
	copy(t.Nodes[pos+1:], t.Nodes[pos:])
	t.Nodes[pos] = node

	return nil
}

func (t Tree) binarySearch(name string) (int, *Node, error) {
	pos := sort.Search(len(t.Nodes), func(i int) bool {
		return t.Nodes[i].Name >= name
	})

	if pos < len(t.Nodes) && t.Nodes[pos].Name == name {
		return pos, t.Nodes[pos], nil
	}

	return pos, nil, errors.New("named node not found")
}

// Subtrees returns a slice of all subtree IDs of the tree.
func (t Tree) Subtrees() (trees IDs) {
	for _, node := range t.Nodes {
		if node.Type == "dir" && node.Subtree != nil {
			trees = append(trees, *node.Subtree)
		}
	}

	return trees
}
