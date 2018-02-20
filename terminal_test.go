package terminal

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

func TestPause(t *testing.T) {
	tests := []struct {
		name string

		message []string

		wantedW string
	}{
		{
			name:    "default message",
			wantedW: "Press enter to continue",
		},
		{
			name:    "custom message",
			message: []string{"do something cool"},
			wantedW: "do something cool",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			readerOriginal := Reader
			writerOriginal := Writer
			Reader = bytes.NewBufferString(test.name + "\n")
			w := bytes.NewBuffer(nil)
			Writer = w
			defer func() {
				Reader = readerOriginal
				Writer = writerOriginal
			}()

			Pause(test.message...)

			if w.String() != test.wantedW {
				t.Errorf("writer contents mismatch: got: %s, wanted: %s", w.String(), test.wantedW)
			}
		})
	}
}

func TestInput(t *testing.T) {
	tests := []struct {
		name string

		reader io.Reader
		delim  byte

		inputName string
		inputVal  []string

		wantedError string
		wantedValue string
		wantedW     string
	}{
		{
			name: "optional",

			reader: bytes.NewBufferString("\n"),

			inputName: "optional",
			inputVal:  []string{Optional},

			wantedValue: "",
			wantedW:     "optional (optional): ",
		},
		{
			name: "optional with default",

			reader: bytes.NewBufferString("\n"),

			inputName: "optional-with-default",
			inputVal:  []string{"default"},

			wantedValue: "default",
			wantedW:     "optional-with-default (optional, default: default): ",
		},
		{
			name: "required",

			reader: bytes.NewBufferString("required value\n"),

			inputName: "required",

			wantedValue: "required value",
			wantedW:     "required (required): ",
		},
		{
			name: "required no input fails",

			reader: bytes.NewBufferString("\n"),

			inputName: "required-no-input",

			wantedError: "required input required-no-input cannot be blank",
			wantedW:     "required-no-input (required): ",
		},
		{
			name: "broken reader",

			reader: new(brokenReader),

			inputName: "broken-reader",
			inputVal:  []string{"default"},

			wantedValue: "default",
			wantedW:     "broken-reader (optional, default: default): ",
		},
		{
			name: "broken reader no default",

			reader: new(brokenReader),

			inputName: "broken-reader-no-default",

			wantedError: "failed to read required input broken-reader-no-default",
			wantedW:     "broken-reader-no-default (required): ",
		},
		{
			name: "custom delim",

			reader: bytes.NewBufferString("te st"),
			delim:  ' ',

			inputName: "custom-delim",

			wantedValue: "te",
			wantedW:     "custom-delim (required): ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			readerOriginal := Reader
			writerOriginal := Writer
			Reader = test.reader
			w := bytes.NewBuffer(nil)
			Writer = w
			delimOriginal := Delim
			if test.delim != 0 {
				Delim = test.delim
			}
			defer func() {
				Reader = readerOriginal
				Writer = writerOriginal
				Delim = delimOriginal
			}()

			gotValue, gotErr := Input(test.inputName, test.inputVal...)
			if errorMessage(gotErr) != test.wantedError {
				t.Errorf("error mismatch: got: %s, wanted: %s", errorMessage(gotErr), test.wantedError)
				return
			}

			if gotValue != test.wantedValue {
				t.Errorf("value mismatch: got: %s, wanted: %s", gotValue, test.wantedValue)
			}

			if w.String() != test.wantedW {
				t.Errorf("writer contents mismatch: got: %s, wanted: %s", w.String(), test.wantedW)
			}
		})
	}
}

func TestMustInput(t *testing.T) {
	tests := []struct {
		name string

		reader io.Reader

		inputName string
		inputVal  []string

		wantedPanic string
		wantedValue string
		wantedW     string
	}{
		{
			name: "optional",

			reader: bytes.NewBufferString("\n"),

			inputName: "optional",
			inputVal:  []string{Optional},

			wantedValue: "",
		},
		{
			name: "broken reader no default",

			reader: new(brokenReader),

			inputName: "broken-reader-no-default",

			wantedPanic: "failed to read required input broken-reader-no-default",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			readerOriginal := Reader
			writerOriginal := Writer
			Reader = test.reader
			Writer = ioutil.Discard
			defer func() {
				Reader = readerOriginal
				Writer = writerOriginal
			}()

			defer func() {
				if r := recover(); r != nil {
					gotPanic := r.(string)
					if r.(string) != test.wantedPanic {
						t.Errorf("panic message mismatch: got: %s, wanted: %s", gotPanic, test.wantedPanic)
					}
				}
			}()

			gotValue := MustInput(test.inputName, test.inputVal...)

			if gotValue != test.wantedValue {
				t.Errorf("value mismatch: got: %s, wanted: %s", gotValue, test.wantedValue)
			}
		})
	}
}

func errorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

type brokenReader struct {
	io.Reader
}

func (br *brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader is broken")
}
