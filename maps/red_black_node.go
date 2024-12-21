package maps

import "github.com/jorge-barroso/collections"

type color bool

const (
	black color = true
	red   color = false
)

// rbNode represents a node in the Red-Black tree
// Kept private to the package to ensure encapsulation
type rbNode[K comparable, V any] struct {
	*collections.Node[Entry[K, V]]
	color  color
	left   *rbNode[K, V]
	right  *rbNode[K, V]
	parent *rbNode[K, V]
}

// newRBNode creates a new Red-Black tree node
func newRBNode[K comparable, V any](entry Entry[K, V]) *rbNode[K, V] {
	return &rbNode[K, V]{
		Node: &collections.Node[Entry[K, V]]{
			Item: entry,
		},
		color: red, // New nodes are always red
	}
}

// isRed returns true if the node is red
func (n *rbNode[K, V]) isRed() bool {
	return n != nil && n.color == red
}

// isBlack returns true if the node is black or nil
func (n *rbNode[K, V]) isBlack() bool {
	return n == nil || n.color == black
}

// getGrandparent returns the grandparent of the node, if it exists
func (n *rbNode[K, V]) getGrandparent() *rbNode[K, V] {
	if n != nil && n.parent != nil {
		return n.parent.parent
	}
	return nil
}

// getUncle returns the uncle of the node, if it exists
func (n *rbNode[K, V]) getUncle() *rbNode[K, V] {
	grandparent := n.getGrandparent()
	if grandparent == nil {
		return nil
	}
	if n.parent == grandparent.left {
		return grandparent.right
	}
	return grandparent.left
}

// getSibling returns the sibling of the node, if it exists
func (n *rbNode[K, V]) getSibling() *rbNode[K, V] {
	if n == nil || n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

// getSuccessor returns the in-order successor of the node
func (n *rbNode[K, V]) getSuccessor() *rbNode[K, V] {
	// If right subtree exists, get minimum of right subtree
	if n.right != nil {
		current := n.right
		for current.left != nil {
			current = current.left
		}
		return current
	}

	// Otherwise, go up until we find a parent where the current node
	// is in the left subtree
	current := n
	parent := current.parent
	for parent != nil && current == parent.right {
		current = parent
		parent = parent.parent
	}
	return parent
}

// getMinimum returns the minimum node in the subtree rooted at this node
func (n *rbNode[K, V]) getMinimum() *rbNode[K, V] {
	current := n
	for current.left != nil {
		current = current.left
	}
	return current
}
