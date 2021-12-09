package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cecobask/kraken-coding-challenge/pkg/transaction"
)

func validateArgs(args []string, wantArgsLen int) (ok bool) {
	if len(args) != wantArgsLen {
		fmt.Fprintf(os.Stderr, "Need to provide %d arguments for the command %s\n", wantArgsLen-1, args[0])
		return
	}
	ok = true
	return
}

func main() {
	stack := transaction.NewStack()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
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
