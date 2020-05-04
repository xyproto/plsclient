package prettypls

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type LanguageServer struct {
	executable            string
	bufIn, bufOut, bufErr bytes.Buffer
	cmd                   *exec.Cmd
	running               bool
	tmpIn                 io.Reader
	tmpOut, tmpErr        io.Writer
}

var errClosed = errors.New("the connection to the language server has been closed")

// New creates a new gopls wrapper
func New() *LanguageServer {
	var ls LanguageServer
	ls.executable = "gopls"
	ls.cmd = exec.Command(ls.executable)
	ls.tmpIn = ls.cmd.Stdin
	ls.cmd.Stdin = &ls.bufIn
	ls.tmpOut = ls.cmd.Stdout
	ls.cmd.Stdout = &ls.bufOut
	ls.tmpOut = ls.cmd.Stderr
	ls.cmd.Stderr = &ls.bufErr
	return &ls
}

// Input will send a string directly to the language server (and start it if needed).
// No headers are added. The resulting stdout and stderr are returned as a string.
func (ls *LanguageServer) Input(s string) (string, error) {
	if ls.cmd == nil {
		return "", errClosed
	}
	if !ls.running {
		ls.bufIn.WriteString(s + "\r")
		if err := ls.cmd.Start(); err != nil {
			return "", err
		}
		ls.running = true
		if err := ls.cmd.Wait(); err != nil {
			output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
		return ls.bufErr.String() + ls.bufOut.String(), nil
	} else {
		ls.bufIn.WriteString(s + "\r")
		if err := ls.cmd.Wait(); err != nil {
			output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
		return ls.bufErr.String() + ls.bufOut.String(), nil
	}
}

// Request will take a JSON message and pass it to the running language server, with the appropriate headers
func (ls *LanguageServer) Request(msg string, verbose bool) (string, error) {
	var sb strings.Builder
	sb.WriteString("Content-Length: " + strconv.Itoa(len(msg)) + "\n")
	sb.WriteString("Content-Type: application/vscode-jsonrpsc; charset=utf-8\n") // Optional
	sb.WriteString("\n")                                                         // Blank line
	sb.WriteString(msg)
	sb.WriteString("\n")
	if verbose {
		fmt.Println(sb.String())
	}
	// The protocol is written by Microsoft, so of course there are DOS line endings
	return ls.Input(strings.Replace(sb.String(), "\n", "\r\n", -1))
}

// Close will close the in/out/err buffers and wait for the process to complete
func (ls *LanguageServer) Close() error {
	if ls.cmd == nil {
		return errClosed
	}
	ls.cmd.Stdin = ls.tmpIn
	ls.cmd.Stdout = ls.tmpOut
	ls.cmd.Stderr = ls.tmpErr
	ls.cmd.Wait()
	ls.cmd = nil
	return nil
}
