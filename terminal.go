package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	// Reader is the reader used to retrieve a value. By default it is os.Stdin.
	Reader io.Reader = os.Stdin
	// Writer is the writer used to which prompts / messages are written. By default it is os.Stdout.
	Writer io.Writer = os.Stdout
	// Delim is the deliminator used to scan for when looking for input. By default it is '\n'.
	Delim = byte('\n')
)

const (
	// Optional is a convenience constant to help define an input as optional.
	Optional = ""
)

// Pause prints an optional message and then returns when a deliminator is read.
func Pause(message ...string) {
	if len(message) > 0 {
		fmt.Fprint(Writer, message[0])
	} else {
		fmt.Fprintf(Writer, "Press enter to continue")
	}
	bufio.NewReader(Reader).ReadString(Delim)
}

// Input prompts for a named string value and returns, with the input value, when a deliminator is read.
// The returned value is automatically trimmed of spaces (including the deliminator if it is a space character).
// If val is specified it is used as a default value for the returned value.
// If val is not specified the value is treated as "required" and an error is returned if no value is entered.
// An error is also returned if there was an error reading from the Reader.
func Input(name string, val ...string) (string, error) {
	if len(val) > 0 {
		if val[0] == Optional {
			fmt.Fprintf(Writer, "%s (optional): ", name)
		} else {
			fmt.Fprintf(Writer, "%s (optional, default: %s): ", name, val[0])
		}
	} else {
		fmt.Fprintf(Writer, "%s (required): ", name)
	}
	reader := bufio.NewReader(Reader)
	str, err := reader.ReadString(Delim)
	if err != nil {
		if len(val) > 0 {
			return val[0], nil
		}
		return "", fmt.Errorf("failed to read required input %s", name)
	}
	str = strings.TrimSpace(str)
	if str == "" {
		if len(val) > 0 {
			return val[0], nil
		}
		return "", fmt.Errorf("required input %s cannot be blank", name)
	}
	return str, nil
}

// MustInput is a convenient form of Input which panics if an error is returned.
func MustInput(name string, val ...string) string {
	s, err := Input(name, val...)
	if err != nil {
		panic(err.Error())
	}
	return s
}
