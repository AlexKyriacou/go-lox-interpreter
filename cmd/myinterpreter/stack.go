package main

import (
	"container/list"
	"sync"
)

type stack[T any] struct {
	dll   *list.List
	mutex sync.Mutex
}

func (s *stack[T]) Len() int {
	return s.dll.Len()
}

func NewStack[T any]() *stack[T] {
	return &stack[T]{dll: list.New()}
}

func (s *stack[T]) Push(x T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.dll.PushBack(x)
}

func (s *stack[T]) Pop() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isEmpty() {
		panic("stack: Pop() called on empty stack")
	}

	tail := s.dll.Back()
	val := tail.Value
	s.dll.Remove(tail)
	return val.(T)
}

func (s *stack[T]) Peek() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isEmpty() {
		panic("stack: Peek() called on empty stack")
	}

	tail := s.dll.Back()
	val := tail.Value
	return val.(T)
}

func (s *stack[T]) isEmpty() bool {
	return s.dll.Len() == 0
}

func (s *stack[T]) get(i int) T{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if i >= s.dll.Len() {
		panic("stack: get() called with index out of bounds")
	}

	e := s.dll.Front()
	for j := 0; j < i; j++ {
		e = e.Next()
	}
	return e.Value.(T)	
}