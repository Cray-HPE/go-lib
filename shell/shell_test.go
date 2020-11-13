package shell

import (
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	sh := Shell{}
	output, err := sh.Exec("echo testing-output", ExecOptions{})
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if output != "testing-output\n" {
		t.Errorf("Expected output 'testing-output\\n', but received '%s'", output)
	}
}

func TestExecTrimmedOutput(t *testing.T) {
	sh := Shell{}
	output, err := sh.Exec("echo testing-output", ExecOptions{TrimOutput: true})
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if output != "testing-output" {
		t.Errorf("Expected output 'testing-output', but received '%s'", output)
	}
}

func TestExecFailure(t *testing.T) {
	sh := Shell{}
	_, err := sh.Exec("command-that-does-not-exist", ExecOptions{})
	if err == nil {
		t.Error("Didn't receive expected failure from command")
	}
}

func TestOutputLineSpliting(t *testing.T) {
	sh := Shell{}
	output, err := sh.Exec(`printf one\ntwo\nthree`, ExecOptions{Silent: true, TrimOutput: true})
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if len(strings.Fields(output)) != 3 {
		t.Errorf("Expected split length of output, expected 3, but received %v", len(strings.Fields(output)))
	}
}

func TestReassembleCommandParts(t *testing.T) {
	sh := Shell{}
	_, err := sh.Exec(`grep -r 'some text with spaces' .`, ExecOptions{Silent: false, TrimOutput: false})
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	_, err = sh.Exec(`grep -r "some text with spaces" .`, ExecOptions{Silent: false, TrimOutput: false})
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
}

func TestGetLines(t *testing.T) {
	multiline := `one
two
three
`
	lines := GetLines(multiline)
	if len(lines) != 3 {
		t.Errorf("Didn't get expected 3 lines from helm.TestGetLines(), instead got: %d", len(lines))
	}
	if lines[0] != "one" || lines[1] != "two" || lines[2] != "three" {
		t.Errorf("Didn't get expected values from split lines in helm.TestGetLines(), instead got: %v", lines)
	}
}
