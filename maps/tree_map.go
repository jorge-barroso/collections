package maps

import (
	"errors"
	"github.com/jorge-barroso/collections"
)

// TreeMap implements a sorted map using a Red-Black tree
type TreeMap[K comparable, V any] struct {
	root *rbNode[K, V]
	size int
	less func(a, b K) bool // Comparison function for keys
}

// NewTreeMap creates a new TreeMap with a custom comparison function
func NewTreeMap[K comparable, V any](less func(a, b K) bool) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		less: less,
	}
}

// Put inserts or updates a key-value pair
func (t *TreeMap[K, V]) Put(key K, value V) {
	entry := Entry[K, V]{
		key:   key,
		value: value,
	}

	newNode := newRBNode(entry)

	if t.root == nil {
		t.root = newNode
		t.size++
		t.insertFixup(newNode)
		return
	}

	t.insert(newNode)
}

// Insert helper methods moved to tree_map_ops.go

// Get retrieves the value associated with a key
func (t *TreeMap[K, V]) Get(key K) (V, error) {
	node := t.findNode(key)
	if node == nil {
		var zero V
		return zero, errors.New("key not found")
	}
	return node.Node.Item.Value(), nil
}

// Remove removes a key-value pair
func (t *TreeMap[K, V]) Remove(key K) error {
	node := t.findNode(key)
	if node == nil {
		return errors.New("key not found")
	}

	t.delete(node)
	t.size--
	return nil
}

// Size returns the number of key-value pairs
func (t *TreeMap[K, V]) Size() int {
	return t.size
}

// NewIterator returns a new iterator for in-order traversal
func (t *TreeMap[K, V]) NewIterator() collections.Iterator[Entry[K, V]] {
	var firstNode *rbNode[K, V]
	if t.root != nil {
		firstNode = t.root.getMinimum()
	}

	return &TreeMapIterator[K, V]{
		tree:    t,
		current: firstNode,
	}
}

// findNode locates a node with the given key
func (t *TreeMap[K, V]) findNode(key K) *rbNode[K, V] {
	current := t.root
	for current != nil {
		if t.less(key, current.Node.Item.Key()) {
			current = current.left
		} else if t.less(current.Node.Item.Key(), key) {
			current = current.right
		} else {
			return current
		}
	}
	return nil
}

func (t *TreeMap[K, V]) insert(node *rbNode[K, V]) {
	// Handle empty tree case first
	if t.root == nil {
		t.root = node
		t.size++
		t.insertFixup(node)
		return
	}

	var parent *rbNode[K, V]
	current := t.root
	// Find insertion point
	for {
		parent = current
		if t.less(node.Node.Item.Key(), current.Node.Item.Key()) {
			current = current.left
		} else if t.less(current.Node.Item.Key(), node.Node.Item.Key()) {
			current = current.right
		} else {
			// Key exists, update value
			current.Node.Item = node.Node.Item
			return
		}

		if current == nil {
			break
		}
	}

	// At this point, parent is guaranteed to be non-nil
	// because we handled the empty tree case earlier
	node.parent = parent
	if t.less(node.Node.Item.Key(), parent.Node.Item.Key()) {
		parent.left = node
	} else {
		parent.right = node
	}

	t.size++
	t.insertFixup(node)
}

func (t *TreeMap[K, V]) insertFixup(node *rbNode[K, V]) {
	for node != t.root && node.parent.isRed() {
		uncle := node.getUncle()
		grandparent := node.getGrandparent()

		if node.parent == grandparent.left {
			if uncle.isRed() {
				node.parent.color = black
				uncle.color = black
				grandparent.color = red
				node = grandparent
			} else {
				if node == node.parent.right {
					node = node.parent
					t.rotateLeft(node)
				}
				node.parent.color = black
				grandparent.color = red
				t.rotateRight(grandparent)
			}
		} else {
			if uncle.isRed() {
				node.parent.color = black
				uncle.color = black
				grandparent.color = red
				node = grandparent
			} else {
				if node == node.parent.left {
					node = node.parent
					t.rotateRight(node)
				}
				node.parent.color = black
				grandparent.color = red
				t.rotateLeft(grandparent)
			}
		}
	}
	t.root.color = black
}

func (t *TreeMap[K, V]) rotateLeft(node *rbNode[K, V]) {
	right := node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.parent = node.parent

	if node.parent == nil {
		t.root = right
	} else if node == node.parent.left {
		node.parent.left = right
	} else {
		node.parent.right = right
	}

	right.left = node
	node.parent = right
}

func (t *TreeMap[K, V]) rotateRight(node *rbNode[K, V]) {
	left := node.left
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.parent = node.parent

	if node.parent == nil {
		t.root = left
	} else if node == node.parent.right {
		node.parent.right = left
	} else {
		node.parent.left = left
	}

	left.right = node
	node.parent = left
}

func (t *TreeMap[K, V]) delete(node *rbNode[K, V]) {
	var successor, child *rbNode[K, V]

	if node.left == nil || node.right == nil {
		successor = node
	} else {
		successor = node.right.getMinimum()
	}

	if successor.left != nil {
		child = successor.left
	} else {
		child = successor.right
	}

	if child != nil {
		child.parent = successor.parent
	}

	if successor.parent == nil {
		t.root = child
	} else if successor == successor.parent.left {
		successor.parent.left = child
	} else {
		successor.parent.right = child
	}

	if successor != node {
		node.Node.Item = successor.Node.Item
	}

	if successor.isBlack() {
		t.deleteFixup(child, successor.parent)
	}
}

func (t *TreeMap[K, V]) deleteFixup(node, parent *rbNode[K, V]) {
	for (node == nil || node.isBlack()) && node != t.root {
		if node == parent.left {
			sibling := parent.right
			if sibling.isRed() {
				sibling.color = black
				parent.color = red
				t.rotateLeft(parent)
				sibling = parent.right
			}
			if sibling.left.isBlack() && sibling.right.isBlack() {
				sibling.color = red
				node = parent
				parent = node.parent
			} else {
				if sibling.right.isBlack() {
					if sibling.left != nil {
						sibling.left.color = black
					}
					sibling.color = red
					t.rotateRight(sibling)
					sibling = parent.right
				}
				sibling.color = parent.color
				parent.color = black
				if sibling.right != nil {
					sibling.right.color = black
				}
				t.rotateLeft(parent)
				node = t.root
			}
		} else {
			sibling := parent.left
			if sibling.isRed() {
				sibling.color = black
				parent.color = red
				t.rotateRight(parent)
				sibling = parent.left
			}
			if sibling.right.isBlack() && sibling.left.isBlack() {
				sibling.color = red
				node = parent
				parent = node.parent
			} else {
				if sibling.left.isBlack() {
					if sibling.right != nil {
						sibling.right.color = black
					}
					sibling.color = red
					t.rotateLeft(sibling)
					sibling = parent.left
				}
				sibling.color = parent.color
				parent.color = black
				if sibling.left != nil {
					sibling.left.color = black
				}
				t.rotateRight(parent)
				node = t.root
			}
		}
	}
	if node != nil {
		node.color = black
	}
}
