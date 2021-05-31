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
  "method": "substract",
  "params": [42, 23],
  "id": 1
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	ls.Close()
	fmt.Println("OUT:\n" + output)
}
