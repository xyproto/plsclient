package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type LanguageServer struct {
	executable            string
	inBuf, outBuf, errBuf bytes.Buffer
	cmd                   *exec.Cmd
	running               bool
}

func (ls *LanguageServer) Input(s string) (string, error) {
	if !ls.running {
		ls.inBuf.WriteString(s + "\n")
		if err := ls.cmd.Start(); err != nil {
			return "", err
		}
		ls.running = true
		if err := ls.cmd.Wait(); err != nil {
			output := strings.TrimSpace(ls.outBuf.String() + ls.errBuf.String())
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
		return ls.outBuf.String() + ls.errBuf.String(), nil
	} else {
		ls.inBuf.WriteString(s + "\n")
		if err := ls.cmd.Wait(); err != nil {
			output := strings.TrimSpace(ls.outBuf.String() + ls.errBuf.String())
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
		return ls.outBuf.String() + ls.errBuf.String(), nil
	}
}

func New() *LanguageServer {
	var ls LanguageServer
	ls.executable = "gopls"
	ls.cmd = exec.Command(ls.executable)
	ls.cmd.Stdin = &ls.inBuf
	ls.cmd.Stdout = &ls.outBuf
	ls.cmd.Stderr = &ls.errBuf
	return &ls
}

func main() {
	ls := New()
	output, err := ls.Input("HELLO")
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
	fmt.Println("OUT:\n" + output)
}
