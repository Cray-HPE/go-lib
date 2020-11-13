// Package shell is a utility for running local shell commands
package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// Interface is the interface for the shell object
type Interface interface {
	Exec(command string, options ExecOptions) (string, error)
}

// Shell is a local command execution utility
type Shell struct{}

// ExecOptions are options for running shell.Exec
type ExecOptions struct {
	Silent     bool
	TrimOutput bool
}

// Exec runs a local shell command
func (sh *Shell) Exec(command string, options ExecOptions) (string, error) {
	command = strings.Replace(command, "\\\n", "", -1)
	commandParts := strings.Split(command, " ")
	reassembledCommandParts := []string{}
	commandPartBuffer := ""
	commandPartBufferQuoteChar := ""
	// the following reconstructs some command line pieces for quoted args with spaces
	// which is a case our simple command line part split doesn't handle above
	// this is very imperfect, but will do for us here until it doesn't
	partSearch, _ := regexp.Compile(`['"]`)
	for _, commandPart := range commandParts {
		matches := partSearch.FindAllStringIndex(commandPart, -1)
		if len(matches) == 1 {
			// we're either starting or ending a buffered command part
			if commandPartBuffer == "" {
				// starting a buffer
				if strings.Contains(commandPart, `"`) {
					commandPartBufferQuoteChar = `"`
				} else {
					commandPartBufferQuoteChar = `'`
				}
				commandPartBuffer = commandPart
			} else if strings.Contains(commandPart, commandPartBufferQuoteChar) {
				// finishing a buffer
				commandPartBuffer = fmt.Sprintf(`%s %s`, commandPartBuffer, commandPart)
				reassembledCommandParts = append(reassembledCommandParts,
					strings.Replace(commandPartBuffer, commandPartBufferQuoteChar, "", -1))
				commandPartBuffer = ""
				commandPartBufferQuoteChar = ""
			}
		} else if commandPartBuffer != "" {
			// in the middle of a re-assemble
			commandPartBuffer = fmt.Sprintf("%s %s", commandPartBuffer, commandPart)
		} else {
			// just a normal re-append
			reassembledCommandParts = append(reassembledCommandParts, commandPart)
		}
	}

	cmd := exec.Command(reassembledCommandParts[0], reassembledCommandParts[1:]...)
	cmd.Env = os.Environ()

	var output bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	cmd.Stdin = os.Stdin
	var stdout io.Writer
	if options.Silent {
		stdout = io.MultiWriter(&output)
	} else {
		stdout = io.MultiWriter(os.Stdout, &output)
	}
	var stderr io.Writer
	if options.Silent {
		stderr = io.MultiWriter(&output)
	} else {
		stderr = io.MultiWriter(os.Stderr, &output)
	}
	err := cmd.Start()
	if err != nil {
		return string(output.Bytes()), err
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return string(output.Bytes()), fmt.Errorf("Shell error: %v", output.String())
	}
	if errStdout != nil {
		return string(output.Bytes()), errStdout
	}
	if errStderr != nil {
		return string(output.Bytes()), errStderr
	}
	if options.TrimOutput {
		return strings.TrimSpace(string(output.Bytes())), nil
	}
	return string(output.Bytes()), nil
}

// GetLines will parse a multiline string and return each line as an item in a slice
func GetLines(multiline string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(multiline))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
