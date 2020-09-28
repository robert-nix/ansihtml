package ansihtml_test

import (
	"testing"

	"github.com/robert-nix/ansihtml"
	"github.com/stretchr/testify/assert"
)

func TestConvertToHTML(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		useClasses  bool
		classPrefix string
		noStyles    bool
		output      string
	}{
		{
			desc:   "no escapes",
			input:  "test",
			output: "test",
		},
		{
			desc:   "html escapes",
			input:  "<test>",
			output: "&lt;test&gt;",
		},
		{
			desc:   "simple color with reset",
			input:  "\x1b[33mYellow\x1b[m",
			output: `<span style="color:olive;">Yellow</span>`,
		},
		{
			desc:   "simple color with no reset",
			input:  "\x1b[33mYellow",
			output: `<span style="color:olive;">Yellow</span>`,
		},
		{
			desc:        "simple color with classnames",
			input:       "\x1b[33mYellow",
			classPrefix: "p-",
			useClasses:  true,
			output:      `<span class="p-fg-yellow">Yellow</span>`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var res []byte
			if tC.useClasses {
				res = ansihtml.ConvertToHTMLWithClasses([]byte(tC.input), tC.classPrefix, tC.noStyles)
			} else {
				res = ansihtml.ConvertToHTML([]byte(tC.input))
			}
			assert.Equal(t, tC.output, string(res), "output")
		})
	}
}
