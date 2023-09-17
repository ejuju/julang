package main

import (
	"os"

	"github.com/ejuju/julang/pkg/julang"
)

func main() {
	// Read code from stdin for REPL mode (default) or from file (if an argument is provided)
	from := os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		from = f
	}

	// Execute code from stdin or file
	vm := julang.NewVM(nil, nil)
	err := vm.Exec(from)
	if err != nil {
		panic(err)
	}
}
