package ioe

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// WORK IKN PROGRESS DOES NOT COMPILE AND IS FULL OF TYPOS ASDFASDFASDFASDF

func Run(cmd, input string) (string,  {
	scdoc := exec.Command("scdoc")

	// Place the current contents in a buffer, and feed it to stdin to the command
	var buf bytes.Buffer
	buf.WriteString(e.String())
	scdoc.Stdin = &buf

	// Create a new file and use it as stdout
	manpageFile, err := os.Create("out.1")
	if err != nil {
		return "", err
	}
	scdoc.Stdout = manpageFile

	var errBuf bytes.Buffer
	scdoc.Stderr = &errBuf

	// Run scdoc
	if err := scdoc.Run(); err != nil {
		statusMessage = strings.TrimSpace(errBuf.String())
		status.ClearAll(c)
		status.SetMessage(statusMessage)
		status.Show(c, e)
		break // from case
	}

	return "", 
}
