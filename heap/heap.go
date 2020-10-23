package myheap

import (
	"container/heap"
)

// Credits: https://golang.org/pkg/container/heap/

// MinHeap is a min heap
type MinHeap []int64

// Len returns the number of items in the MinHeap
func (h MinHeap) Len() int { return len(h) }

// Less returns true if the MinHeap[i] < MinHeap[j]
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }

// Swap swaps elements in the MinHeap
func (h MinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push an int64 to MinHeap
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

// Pop minimum in64 in the MinHeap
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MaxHeap is a max heap
type MaxHeap []int64

// Len returns the number of items in the MaxHeap
func (h MaxHeap) Len() int { return len(h) }

// Less returns true if the MaxHeap[i] > MaxHeap[j]
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }

// Swap swaps elements in the MaxHeap
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push an int64 to MaxHeap
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

// Pop maximum in64 in the MaxHeap
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Test is a dummy test
func Test() {
	maxHeap := &MaxHeap{1, 2, 3}
	heap.Init(maxHeap)
}
