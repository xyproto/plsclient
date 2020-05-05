package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xyproto/prettypls"
)

func main() {
	ls := prettypls.New()
	defer ls.Close()

	curdir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	_ = curdir

	output, err := ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": `+strconv.Itoa(os.Getpid())+`,
    "rootPath": "`+curdir+`",
    "rootUri": null,
    "capabilities": {
    }
  }
}`, true)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// return
	}
	fmt.Println("OUT:\n" + output)
}
