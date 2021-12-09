package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type writeAction struct {
	value string
}

type deleteAction struct{}

type Transaction struct {
	parent  *Transaction
	actions map[string]interface{}
}

type Stack struct {
	top   *Transaction
	store map[string]string
}

func NewTransaction(previousTransaction *Transaction) Transaction {
	return Transaction{
		parent:  previousTransaction,
		actions: make(map[string]interface{}),
	}
}

func NewStack() Stack {
	return Stack{
		store: make(map[string]string),
	}
}

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
		fmt.Fprintln(os.Stderr, "No current transaction, COMMIT is not possible")
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

func validateArgs(args []string, wantArgsLen int) (ok bool) {
	if len(args) != wantArgsLen {
		fmt.Fprintf(os.Stderr, "Need to provide %d arguments for the command %s\n", wantArgsLen-1, args[0])
		return
	}
	ok = true
	return
}

func main() {
	stack := NewStack()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		args[0] = strings.ToUpper(args[0])
		switch args[0] {
		case "READ":
			if ok := validateArgs(args, 2); !ok {
				continue
			}
			stack.Read(args[1])
		case "WRITE":
			if ok := validateArgs(args, 3); !ok {
				continue
			}
			stack.Write(args[1], args[2])
		case "DELETE":
			validateArgs(args, 2)
			if ok := validateArgs(args, 2); !ok {
				continue
			}
			stack.Delete(args[1])
		case "START":
			if ok := validateArgs(args, 1); !ok {
				continue
			}
			stack.Start()
		case "COMMIT":
			if ok := validateArgs(args, 1); !ok {
				continue
			}
			stack.Commit()
		case "ABORT":
			if ok := validateArgs(args, 1); !ok {
				continue
			}
			stack.Abort()
		case "QUIT":
			if ok := validateArgs(args, 1); !ok {
				continue
			}
			fmt.Fprintln(os.Stdout, "Exiting...")
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "Invalid command")
		}
	}
}
