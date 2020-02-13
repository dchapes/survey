package survey

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestMultilineRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Multiline
		data     MultilineTemplateData
		expected string
	}{
		{
			"Test Multiline question output without default",
			Multiline{Message: "What is your favorite month:"},
			MultilineTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", defaultIcons().Question.Text),
		},
		{
			"Test Multiline question output with default",
			Multiline{Message: "What is your favorite month:", Default: "April"},
			MultilineTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", defaultIcons().Question.Text),
		},
		{
			"Test Multiline answer output",
			Multiline{Message: "What is your favorite month:"},
			MultilineTemplateData{Answer: "October", ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: \nOctober", defaultIcons().Question.Text),
		},
		{
			"Test Multiline question output without default but with help hidden",
			Multiline{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter 2 empty lines to finish]", defaultPromptConfig().HelpInput),
		},
		{
			"Test Multiline question output with default and with help hidden",
			Multiline{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter 2 empty lines to finish]", defaultPromptConfig().HelpInput),
		},
		{
			"Test Multiline question output without default but with help shown",
			Multiline{Message: "What is your favorite month:", Help: "This is helpful"},
			MultilineTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: [Enter 2 empty lines to finish]", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Multiline question output with default and with help shown",
			Multiline{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			MultilineTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) [Enter 2 empty lines to finish]", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
	}

	var sb strings.Builder
	for _, test := range tests {
		sb.Reset()
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal("os.Pipe:", err)
		}

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Multiline = test.prompt
		// set the icon set
		test.data.Config = defaultPromptConfig()

		err = test.prompt.Render(
			MultilineQuestionTemplate,
			test.data,
		)
		if err != nil {
			t.Errorf("%s test.prompt.Render() failed:\n\t%v", test.title, err)
			//continue
		}

		if err := w.Close(); err != nil {
			t.Fatal("os.Pipe w.Close():", err)
		}
		if _, err := io.Copy(&sb, r); err != nil {
			t.Fatal("io.Copy():", err)
		}

		if g, w := sb.String(), test.expected; !strings.Contains(g, w) {
			t.Errorf("%s\ngave %q\n\twanted to contain %q", test.title, g, w)
		}
	}
}

func TestMultilinePrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Multiline prompt interaction",
			&Multiline{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird\nI guess...\nnot sure\n\n")
				c.ExpectEOF()
			},
			"Larry Bird\nI guess...\nnot sure",
		},
		{
			"Test Multiline prompt interaction with default",
			&Multiline{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("\n\n")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test Multiline prompt interaction overriding default",
			&Multiline{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird\n\n")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Multiline does not implement help interaction",
			&Multiline{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.SendLine("Satoshi Nakamoto\n\n")
				c.ExpectEOF()
			},
			"?\nSatoshi Nakamoto",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
