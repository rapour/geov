package geov

import (
	"golang.org/x/exp/constraints"
)

func NewHeap[T constraints.Ordered](arr []Element[T]) Heap[T] {
	return Heap[T]{array: arr}
}

type Element[T constraints.Ordered] interface {
	Value() T
}

type Heap[T constraints.Ordered] struct {
	array []Element[T]
}

func (h *Heap[T]) Add(e Element[T]) {
	h.array = append([]Element[T]{e}, h.array...)
	h.MinHeapify(0)
}

func (h *Heap[T]) Min() Element[T] { return h.array[0] }

func (h *Heap[T]) ExtractMin() Element[T] {

	if h.GetSize() == 0 {
		return nil
	}

	min := h.array[0]
	last := h.array[h.GetSize()-1]

	h.array = h.array[:h.GetSize()-1]

	if min == last {
		return min
	}

	h.array[0] = last
	h.MinHeapify(0)

	return min
}

func (h *Heap[T]) GetSize() int {
	return len(h.array)
}

func (h *Heap[T]) Left(i int) int {
	return 2*i + 1
}

func (h *Heap[T]) Right(i int) int {
	return 2*i + 2
}

func (h *Heap[T]) MinHeapify(i int) {

	left := h.Left(i)
	right := h.Right(i)

	smallest := h.array[i].Value()
	smallestIndex := i

	if left < h.GetSize() && h.array[left].Value() < smallest {
		smallest = h.array[left].Value()
		smallestIndex = left
	}

	if right < h.GetSize() && h.array[right].Value() < smallest {
		smallest = h.array[right].Value()
		smallestIndex = right
	}

	if smallestIndex != i {
		h.array[i], h.array[smallestIndex] = h.array[smallestIndex], h.array[i]
		h.MinHeapify(smallestIndex)
	}

}

func (h *Heap[T]) BuildMinHeap() {
	lastNodeIndx := h.GetSize() - 1

	for i := lastNodeIndx; i >= 0; i-- {
		h.MinHeapify(i)
	}
}
