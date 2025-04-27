package main

import (
	"encoding/hex"
	"fmt"
	"github.com/minio/sha256-simd"
	"log"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	printList(list1)
	printList(list2)
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

func printList(list *ListNode) {
	for list != nil {
		fmt.Print(list.Val, "\t")
		list = list.Next
	}
	fmt.Println()
}
func main() {
	list1 := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 4}}}
	list2 := &ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}}

	result := mergeTwoLists(list1, list2)
	printList(result)

	hashedKey := RoutingKeyCalculator(
		"2061fcaf013131a753bac07e10cdf46eae95cb96bbbfcdbd7564667fc350db62")
	log.Println(hashedKey)
}

func RoutingKeyCalculator(key string) string {
	if len(key) != 64 {
		log.Println("Ivalid Validator Key")
		return ""
	}

	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
