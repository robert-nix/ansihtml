package ansihtml_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/robert-nix/ansihtml"
	"github.com/stretchr/testify/assert"
)

type escapeParams struct {
	finalByte         byte
	intermediateBytes []byte
	parameterBytes    []byte
}

func TestParser(t *testing.T) {
	testCases := []struct {
		desc       string
		input      string
		output     string
		escapes    []escapeParams
		bufferSize int
	}{
		{
			desc:    "no escapes",
			input:   "test",
			output:  "test",
			escapes: nil,
		},
		{
			desc:   "only set color",
			input:  "\x1b[0;33m",
			output: "",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
			},
		},
		{
			desc:   "term title",
			input:  "\x1b]0;window title\x1b\\",
			output: "",
			escapes: []escapeParams{
				{
					finalByte:         ']',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;window title\x1b\\"),
				},
			},
		},
		{
			desc:   "xterm title",
			input:  "\x1b]0;window title\x07",
			output: "",
			escapes: []escapeParams{
				{
					finalByte:         ']',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;window title\x07"),
				},
			},
		},
		{
			desc:   "unicode escapes",
			input:  "\u009b0;33m",
			output: "",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
			},
		},
		{
			desc:    "unicode almost escape",
			input:   "test\u00a00;33m",
			output:  "test\u00a00;33m",
			escapes: nil,
		},
		{
			desc:   "unicode escapes split across buffer boundary",
			input:  "test\u009b0;33m",
			output: "test",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
			},
			bufferSize: 5,
		},
		{
			desc:       "unicode almost escape split across buffer boundary",
			input:      "test\u00a00;33m",
			output:     "test\u00a00;33m",
			escapes:    nil,
			bufferSize: 5,
		},
		{
			desc:   "multiple escapes",
			input:  "test\x1b[0;33mtest\x1b[m",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("m"),
				},
			},
		},
		{
			desc:   "multiple escapes in a row",
			input:  "test\x1b[0;33m\x1b[mtest",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("m"),
				},
			},
		},
		{
			desc:   "multiple escapes with buffer length 1",
			input:  "test\x1b[0;33mtest\x1b[m",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("0;33m"),
				},
				{
					finalByte:         '[',
					intermediateBytes: nil,
					parameterBytes:    []byte("m"),
				},
			},
			bufferSize: 1,
		},
		{
			desc:   "paramless escape",
			input:  "test\x1bctest",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         'c',
					intermediateBytes: nil,
					parameterBytes:    nil,
				},
			},
		},
		{
			desc:   "paramless unicode escape",
			input:  "test\u009ctest",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         '\\',
					intermediateBytes: nil,
					parameterBytes:    nil,
				},
			},
		},
		{
			desc:    "unknown sequence",
			input:   "test\x1b\x7ftest",
			output:  "testtest",
			escapes: nil,
		},
		{
			desc:    "invalid CSI first byte",
			input:   "test\x1b[\x7ftest",
			output:  "testtest",
			escapes: nil,
		},
		{
			desc:    "invalid CSI second byte",
			input:   "test\x1b[0 \x7ftest",
			output:  "testtest",
			escapes: nil,
		},
		{
			desc:   "intermediate byte",
			input:  "test\x1b Ftest",
			output: "testtest",
			escapes: []escapeParams{
				{
					finalByte:         'F',
					intermediateBytes: []byte(" "),
					parameterBytes:    nil,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rd := strings.NewReader(tC.input)
			w := new(strings.Builder)
			var ei int
			p := ansihtml.NewParser(rd, w)
			handler := func(finalByte byte, intermediateBytes []byte, parameterBytes []byte) {
				if !assert.Less(t, ei, len(tC.escapes), "too many escapes") {
					return
				}
				e := tC.escapes[ei]
				ei++
				assert.Equal(t, e.finalByte, finalByte, "finalByte")
				assert.Equal(t, e.intermediateBytes, intermediateBytes, "intermediateBytes")
				assert.Equal(t, e.parameterBytes, parameterBytes, "parameterBytes")
			}
			var err error
			if tC.bufferSize == 0 {
				err = p.Parse(handler)
			} else {
				err = p.ParseBuffer(make([]byte, tC.bufferSize), handler)
			}
			assert.Equal(t, err, nil)
			assert.Equal(t, ei, len(tC.escapes), "too few escapes")
			assert.Equal(t, tC.output, w.String(), "output")
		})
	}
}

func TestParseBuffer(t *testing.T) {
	p := ansihtml.NewParser(nil, nil)
	err := p.ParseBuffer(nil, nil)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "buffer must not be empty")
	}
}

type errorWriter struct{}

const errorWriterErr = "cannot write to errorWriter"

func (w *errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New(errorWriterErr)
}

func TestParseEmptyWriter(t *testing.T) {
	rd := strings.NewReader("this is \x1b[0;33myellow\x1b[m")
	var w errorWriter
	p := ansihtml.NewParser(rd, &w)
	err := p.Parse(nil)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), errorWriterErr)
	}
}

type errorReader struct{}

const errorReaderErr = "cannot read from errorReader"

func (r *errorReader) Read(b []byte) (int, error) {
	return 0, errors.New(errorReaderErr)
}

func TestParseErrorReader(t *testing.T) {
	rd := errorReader{}
	w := strings.Builder{}
	p := ansihtml.NewParser(&rd, &w)
	err := p.Parse(nil)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), errorReaderErr)
	}
}
