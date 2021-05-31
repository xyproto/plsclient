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

const dosnl = "\r\n"
const nl = "\n"
const r = "\r"

// LanguageServer wraps a language server executable and keeps track of its stdin/stdout/stderr
type LanguageServer struct {
	Running               bool
	bufIn, bufOut, bufErr bytes.Buffer
	cmd                   *exec.Cmd
	tmpIn                 io.Reader
	tmpOut, tmpErr        io.Writer
}

var errClosed = errors.New("the connection to the language server has been closed")

// NewCustom creates a new language server wrapper that wraps a custom command.
func NewCustom(executable string) *LanguageServer {
	var ls LanguageServer
	ls.cmd = exec.Command(executable)
	ls.tmpIn = ls.cmd.Stdin
	ls.cmd.Stdin = &ls.bufIn
	ls.tmpOut = ls.cmd.Stdout
	ls.cmd.Stdout = &ls.bufOut
	ls.tmpOut = ls.cmd.Stderr
	ls.cmd.Stderr = &ls.bufErr
	return &ls
}

// New creates a new LanguageServer that wraps the "gopls" command.
func New() *LanguageServer {
	return NewCustom("gopls")
}

// SendIn will send a string directly to the language server stdin (and start it if needed).
// No headers are added. The resulting stdout and stderr are returned as a string.
func (ls *LanguageServer) SendIn(s string, verbose bool) (string, error) {
	if ls.cmd == nil {
		return "", errClosed
	}
	if !ls.Running {
		if verbose {
			fmt.Printf("Writing string of length %d to language server stdin\n", len(s))
		}
		ls.bufIn.WriteString(s + dosnl)
		if err := ls.cmd.Start(); err != nil {
			return "", err
		}
		ls.Running = true
		if err := ls.cmd.Wait(); err != nil {
			//if verbose {
				//fmt.Println("1 GOT ERROR:", err)
			//}
			output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}

		if verbose {
			fmt.Println("1 GOT STDERR")
			fmt.Println(ls.bufErr.String())
			fmt.Println("1 GOT STDOUT")
			fmt.Println(ls.bufOut.String())
		}

		return ls.bufErr.String() + ls.bufOut.String(), nil
	}
	ls.bufIn.WriteString(s + dosnl)
	if err := ls.cmd.Wait(); err != nil {
		if verbose {
			fmt.Println("2 GOT ERROR:", err)
		}
		output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
		if output != "" {
			return "", errors.New(output)
		}
		return "", err
	}
	if verbose {
		fmt.Println("2 GOT STDERR")
		fmt.Println(ls.bufErr.String())
		fmt.Println("2 GOT STDOUT")
		fmt.Println(ls.bufOut.String())
	}

	return ls.bufErr.String() + ls.bufOut.String(), nil
}

// SendInBytes will send bytes directly to the language server stdin (and start it if needed).
// No headers are added. The resulting stdout and stderr are returned as a byte slice.
func (ls *LanguageServer) SendInBytes(b []byte, verbose bool) ([]byte, error) {
	if ls.cmd == nil {
		return []byte{}, errClosed
	}
	if !ls.Running {
		if verbose {
			fmt.Printf("Writing bytes of length %d to language server stdin\n", len(b))
		}
		ls.bufIn.Write(b)
		ls.bufIn.WriteString(dosnl)
		if err := ls.cmd.Start(); err != nil {
			return []byte{}, err
		}
		ls.Running = true
		if err := ls.cmd.Wait(); err != nil {
			output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
			if output != "" {
				return []byte{}, errors.New(output)
			}
			return []byte{}, err
		}
		var bufCombined bytes.Buffer
		bufCombined.Write(ls.bufErr.Bytes())
		bufCombined.Write(ls.bufOut.Bytes())
		return bufCombined.Bytes(), nil
	}
	ls.bufIn.Write(b)
	ls.bufIn.WriteString(dosnl)
	if err := ls.cmd.Wait(); err != nil {
		output := strings.TrimSpace(ls.bufOut.String() + ls.bufErr.String())
		if output != "" {
			return []byte{}, errors.New(output)
		}
		return []byte{}, err
	}
	var bufCombined bytes.Buffer
	bufCombined.Write(ls.bufErr.Bytes())
	bufCombined.Write(ls.bufOut.Bytes())
	return bufCombined.Bytes(), nil
}

// Request will take a JSON string and pass it to the running language server, with the appropriate headers
func (ls *LanguageServer) Request(msg string, verbose bool) (string, error) {
	// The protocol is written by Microsoft, so of course there are DOS line endings in the JSON data
	dosMessage := msg
	if !strings.Contains(msg, dosnl) {
		dosMessage = strings.ReplaceAll(msg, nl, dosnl)
	}

	// Build the request string
	var sb strings.Builder
	sb.WriteString("Content-Length: " + strconv.Itoa(len(dosMessage)) + r)
	//sb.WriteString("Content-Type: application/vscode-jsonrpsc; charset=utf-8" + dosnl) // Optional
	sb.WriteString(nl + r)                                                          // Blank line
	sb.WriteString(dosMessage)
	if verbose {
		fmt.Println(strings.ReplaceAll(sb.String(), dosnl, nl))
	}
	return ls.SendIn(sb.String(), verbose)
}

// RequestBytes will take JSON bytes and pass it to the running language server, with the appropriate headers
func (ls *LanguageServer) RequestBytes(msg []byte, verbose bool) ([]byte, error) {
	// The protocol is written by Microsoft, so of course there are DOS line endings in the JSON data
	dosMessage := msg
	if !bytes.Contains(msg, []byte(dosnl)) {
		dosMessage = bytes.ReplaceAll(msg, []byte(nl), []byte(dosnl))
	}

	// Build the request byte array
	var buf bytes.Buffer
	buf.WriteString("Content-Length: " + strconv.Itoa(len(dosMessage)) + r)
	//buf.WriteString("Content-Type: application/vscode-jsonrpsc; charset=utf-8" + dosnl) // Optional
	buf.WriteString(nl + r)                                                             // Blank line
	buf.Write(dosMessage)
	if verbose {
		fmt.Println(strings.ReplaceAll(buf.String(), dosnl, nl))
	}
	return ls.SendInBytes(buf.Bytes(), verbose)
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
	ls.Running = false
	return nil
}
