package main_test

import (
	"fmt"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

//	func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
//		for l1 != nil && l2 != nil {
//			return nil
//		}
//		leftSide := 0
//		for l1 != nil {
//			leftSide += l1.Val
//		}
//		rightSide := 0
//		for l2 != nil {
//			rightSide += l2.Val
//		}
//		sum := leftSide + rightSide
//
// }
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummyHead := &ListNode{}
	current := dummyHead

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			current.Next = list1
			list1 = list1.Next
		} else {
			current.Next = list2
			list2 = list2.Next
		}
		current = current.Next
	}
	if list1 != nil {
		current.Next = list1
	} else {
		current.Next = list2
	}
	return current.Next
}
func mergeTwoLists2(list1 *ListNode, list2 *ListNode) *ListNode {
	fmt.Println(list1, list2)
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			list1.Next = mergeTwoLists(list1.Next, list2)
			return list1
		} else {
			list2.Next = mergeTwoLists(list2.Next, list1)
			return list2
		}

	}
	if list1 == nil {
		return list2
	}
	return list1
}
func printList(head *ListNode) {
	for head != nil {
		fmt.Print(head.Val, " ")
		head = head.Next
	}
	fmt.Println()
}

func TestMergeTwoLists(t *testing.T) {
	// Creating test case 1: [1, 2, 4] + [1, 3, 4]
	list1 := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 4}}}
	list2 := &ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}}

	// Expected merged list: [1, 1, 2, 3, 4, 4]
	result := mergeTwoLists2(list1, list2)

	// Print the result for manual inspection
	fmt.Print("Merged List: ")
	printList(result)

	// Optionally, you can convert print to assertions for automated testing
	expected := []int{1, 1, 2, 3, 4, 4}
	var actual []int
	for result != nil {
		actual = append(actual, result.Val)
		result = result.Next
	}

	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Test failed: expected %v, got %v", expected, actual)
			return
		}
	}
}
