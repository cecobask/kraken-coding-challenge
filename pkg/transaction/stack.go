package transaction

import (
	"fmt"
	"os"
)

type Stack struct {
	top   *Transaction
	store map[string]string
}

func NewStack() Stack {
	return Stack{
		store: make(map[string]string),
	}
}

type writeAction struct {
	value string
}

type deleteAction struct{}

func (s *Stack) Write(key string, value string) {
	if s.top != nil {
		s.top.actions[key] = writeAction{value: value}
	} else {
		s.store[key] = value
	}
}

func (s *Stack) Delete(key string) {
	if s.top != nil {
		s.top.actions[key] = deleteAction{}
	} else {
		delete(s.store, key)
	}
}

func (s *Stack) Read(key string) {
	for s.top != nil {
		if _, ok := s.top.actions[key]; !ok {
			s.top = s.top.parent
		} else {
			switch s.top.actions[key].(type) {
			case deleteAction:
				fmt.Fprintln(os.Stderr, "Key not found:", key)
				return
			case writeAction:
				fmt.Fprintln(os.Stdout, s.top.actions[key].(writeAction).value)
				return
			}
		}
	}
	if value, ok := s.store[key]; ok {
		fmt.Fprintln(os.Stdout, value)
		return
	}
	fmt.Fprintln(os.Stderr, "Key not found:", key)
}

func (s *Stack) Start() {
	transaction := NewTransaction(s.top)
	s.top = &transaction
}

func (s *Stack) Abort() {
	if s.top == nil {
		fmt.Fprintln(os.Stderr, "No current transaction, ABORT is not possible")
		return
	}
	s.top = s.top.parent
}

func (s *Stack) Commit() {
	if s.top == nil {
		return
	}
	if s.top.parent == nil {
		for key, value := range s.top.actions {
			switch value.(type) {
			case writeAction:
				s.store[key] = value.(writeAction).value
			case deleteAction:
				delete(s.store, key)
			}
			s.top = nil
		}
	} else {
		for key, value := range s.top.actions {
			s.top.parent.actions[key] = value
		}
		s.top = s.top.parent
	}
}
