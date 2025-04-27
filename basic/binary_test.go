package main_test

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"testing"
)

type Node[T constraints.Ordered] struct {
	Val   T
	Left  *Node[T]
	Right *Node[T]
}

func (n *Node[T]) Insert(val T) {
	if n.Val == val {
		return
	}
	if val < n.Val {
		if n.Left == nil {
			n.Left = &Node[T]{Val: val}
			return
		} else {
			n.Left.Insert(val)
		}
	} else {
		if n.Right == nil {
			n.Right = &Node[T]{Val: val}
			return
		} else {
			n.Right.Insert(val)
		}
	}
}

func (n *Node[T]) InOrder() {
	if n == nil {
		return
	}
	n.Left.InOrder()
	fmt.Println(n.Val)
	n.Right.InOrder()
}
func (n *Node[T]) Search(val T) bool {
	if n.Val == val {
		return true
	}
	if val < n.Val {
		if n.Left == nil {
			return false
		}
		return n.Left.Search(val)
	} else {
		if n.Right == nil {
			return false
		}
		return n.Right.Search(val)
	}
}

func buildIntTree() *Node[int] {
	root := &Node[int]{Val: 10}
	root.Insert(5)
	root.Insert(15)
	root.Insert(3)
	root.Insert(7)
	root.Insert(12)
	root.Insert(17)
	return root
}

func TestSearchIntTree(t *testing.T) {
	tree := buildIntTree()

	tests := []struct {
		input    int
		expected bool
	}{
		{5, true},
		{12, true},
		{17, true},
		{8, false},
		{0, false},
		{100, false},
	}

	for _, test := range tests {
		result := tree.Search(test.input)
		if result != test.expected {
			t.Errorf("Search(%d) = %v; want %v", test.input, result, test.expected)
		}
	}
}
