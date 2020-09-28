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
		{
			desc:  "every style",
			input: "\x1b[1;3;4;5;7;8;9;11;26;38;5;255;51;53;58;2;0;255;0mMany styles\x1b[2;20;21;6;12;38;2;233;277;255mA few more styles",
			output: `<span style="font-weight:bold;font-style:italic;text-decoration-line:underline line-through overline;` +
				`filter:invert(100%);opacity:0;font-family:sans-serif;color:rgb(238,238,238);text-decoration-color:rgb(0,255,0);">` +
				`Many styles</span><span style="font-weight:lighter;text-decoration-line:line-through overline;filter:invert(100%);` +
				`opacity:0;font-family:sans-serif;color:rgb(233,21,255);text-decoration-color:rgb(0,255,0);">A few more styles</span>`,
		},
		{
			desc:       "every class",
			input:      "\x1b[1;3;4;5;7;8;9;11;26;38;5;255;51;53;58;2;0;255;mMany classes\x1b[2;20;21;6;12;38;2;233;277;255mA few more classes",
			useClasses: true,
			output: `<span class="bold italic underline strikethrough overline slow-blink invert hide font-1 proportional" style="` +
				`color:rgb(238,238,238);text-decoration-color:rgb(0,255,0);">Many classes</span><span class="faint fraktur ` +
				`double-underline strikethrough overline fast-blink invert hide font-2 proportional" style="color:rgb(233,21,255);` +
				`text-decoration-color:rgb(0,255,0);">A few more classes</span>`,
		},
		{
			desc:       "every class no styles",
			input:      "\x1b[1;3;4;5;7;8;9;11;26;38;5;255;51;53;58;2;0;255;mMany classes\x1b[2;20;21;6;12;38;2;233;277;255mA few more classes",
			useClasses: true,
			noStyles:   true,
			output: `<span class="bold italic underline strikethrough overline slow-blink invert hide font-1 proportional">Many classes` +
				`</span><span class="faint fraktur double-underline strikethrough overline fast-blink invert hide font-2 proportional">` +
				`A few more classes</span>`,
		},
		{
			desc:       "background color",
			input:      "\x1b[48;5;28mGreen background\x1b[102mDifferent green background",
			useClasses: true,
			output: `<span style="background-color:rgb(0,135,0);">Green background</span>` +
				`<span class="bg-bright-green">Different green background</span>`,
		},
		{
			desc:       "more colors",
			input:      "\x1b[47;30mBlack on white\x1b[97;40mWhite on black",
			useClasses: true,
			output: `<span class="fg-black bg-white">Black on white</span>` +
				`<span class="fg-bright-white bg-black">White on black</span>`,
		},
		{
			desc:  "superscript and subscript",
			input: "\x1b[73mSuperscript\x1b[74mSubscript",
			output: `<span style="vertical-align:super;">Superscript</span>` +
				`<span style="vertical-align:sub;">Subscript</span>`,
		},
		{
			desc:       "superscript and subscript classes",
			input:      "\x1b[73mSuperscript\x1b[74mSubscript",
			useClasses: true,
			output: `<span class="superscript">Superscript</span>` +
				`<span class="subscript">Subscript</span>`,
		},
		{
			desc:   "integer overflow saturates",
			input:  "\x1b[30;48;2;9999999997;9999999998;9999999999mThis is black on white",
			output: `<span style="color:black;background-color:rgb(255,255,255);">This is black on white</span>`,
		},
		{
			desc:   "noops",
			input:  "\x1bc\x1b[2J\x1b[?1m\x1b[38;1mtest",
			output: `test`,
		},
		{
			desc:       "resets",
			input:      "\x1b[1;3;4;5;7;8;9;31;47;26;52;53;58;5;2mtest\x1b[59;55;54;50;49;39;29;28;27;25;24;23;22mtest",
			useClasses: true,
			output:     `<span class="bold italic underline strikethrough overline slow-blink invert hide proportional fg-red bg-white underline-green">test</span>test`,
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
