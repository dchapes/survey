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

func TestConfirmRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Confirm
		data     ConfirmTemplateData
		expected string
	}{
		{
			"Test Confirm question output with default true",
			Confirm{Message: "Is pizza your favorite food?", Default: true},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? (Y/n) ", defaultIcons().Question.Text),
		},
		{
			"Test Confirm question output with default false",
			Confirm{Message: "Is pizza your favorite food?", Default: false},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? (y/N) ", defaultIcons().Question.Text),
		},
		{
			"Test Confirm answer output",
			Confirm{Message: "Is pizza your favorite food?"},
			ConfirmTemplateData{Answer: "Yes"},
			fmt.Sprintf("%s Is pizza your favorite food? Yes\n", defaultIcons().Question.Text),
		},
		{
			"Test Confirm with help but help message is hidden",
			Confirm{Message: "Is pizza your favorite food?", Help: "This is helpful"},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? [%s for help] (y/N) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Confirm help output with help message shown",
			Confirm{Message: "Is pizza your favorite food?", Help: "This is helpful"},
			ConfirmTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s Is pizza your favorite food? (y/N) ", defaultIcons().Help.Text, defaultIcons().Question.Text),
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
		test.data.Confirm = test.prompt

		// set the runtime config
		test.data.Config = defaultPromptConfig()

		err = test.prompt.Render(
			ConfirmQuestionTemplate,
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

func TestConfirmPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Confirm prompt interaction",
			&Confirm{
				Message: "Is pizza your favorite food?",
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (y/N)")
				c.SendLine("n")
				c.ExpectEOF()
			},
			false,
		},
		{
			"Test Confirm prompt interaction with default",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Default: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (Y/n)")
				c.SendLine("")
				c.ExpectEOF()
			},
			true,
		},
		{
			"Test Confirm prompt interaction overriding default",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Default: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (Y/n)")
				c.SendLine("n")
				c.ExpectEOF()
			},
			false,
		},
		{
			"Test Confirm prompt interaction and prompt for help",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Help:    "It probably is",
			},
			func(c *expect.Console) {
				c.ExpectString(
					fmt.Sprintf(
						"Is pizza your favorite food? [%s for help] (y/N)",
						string(defaultPromptConfig().HelpInput),
					),
				)
				c.SendLine("?")
				c.ExpectString("It probably is")
				c.SendLine("Y")
				c.ExpectEOF()
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
