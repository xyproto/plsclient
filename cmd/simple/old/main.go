package main

import (
	"fmt"
	"os"

	"github.com/xyproto/prettypls"
)

func main() {
	ls := prettypls.New()

	output, err := ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "substract",
  "params": [42, 23]
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	fmt.Println("OUT:\n" + output)
	//ls.Close()

	output, err = ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "textDocument/didOpen",
  "params": {}
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	fmt.Println("OUT:\n" + output)
	//ls.Close()

	output, err = ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "textDocument/didOpen",
  "params": {}
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	fmt.Println("OUT:\n" + output)
	//ls.Close()

	output, err = ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "textDocument/didOpen",
  "params": {}
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	fmt.Println("OUT:\n" + output)

	ls.Close()
}
