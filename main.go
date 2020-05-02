package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type LanguageServer struct {
	executable            string
	inBuf, outBuf, errBuf bytes.Buffer
	cmd                   *exec.Cmd
	running               bool
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

func (ls *LanguageServer) Input(s string) (string, error) {
	if !ls.running {
		ls.inBuf.WriteString(s + "\r")
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
		ls.inBuf.WriteString(s + "\r")
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

func (ls *LanguageServer) Request(msg string) (string, error) {
	var sb strings.Builder
	sb.WriteString("Content-Length:" + strconv.Itoa(len(msg)) + "\n")
	//sb.WriteString("Content-Type: application/vscode-jsonrpsc;charset=utf-8\n")
	sb.WriteString("\n")
	sb.WriteString(msg)
	sb.WriteString("\n")
	fmt.Println("REQUEST:\n" + sb.String())
	return ls.Input(strings.Replace(sb.String(), "\n", "\r\n", -1))
}

func main() {
	ls := New()
	output, err := ls.Request(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "textDocument/didOpen",
  "params": {}
}`)
	if err != nil {
		if output != "" {
			fmt.Printf("output: %s\n", output)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
	fmt.Println("OUT:\n" + output)
}
