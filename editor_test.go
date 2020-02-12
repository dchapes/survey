package survey

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestEditorRender(t *testing.T) {
	tests := []struct {
		title    string
		prompt   Editor
		data     EditorTemplateData
		expected string
	}{
		{
			"Test Editor question output without default",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with default",
			Editor{Message: "What is your favorite month:", Default: "April"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with HideDefault",
			Editor{Message: "What is your favorite month:", Default: "April", HideDefault: true},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Question.Text),
		},
		{
			"Test Editor answer output",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{Answer: "October", ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: October\n", defaultIcons().Question.Text),
		},
		{
			"Test Editor question output without default but with help hidden",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] [Enter to launch editor] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Editor question output with default and with help hidden",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] (April) [Enter to launch editor] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Editor question output without default but with help shown",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: [Enter to launch editor] ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Editor question output with default and with help shown",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) [Enter to launch editor] ", defaultIcons().Help.Text, defaultIcons().Question.Text),
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
		test.data.Editor = test.prompt

		// set the icon set
		test.data.Config = defaultPromptConfig()

		err = test.prompt.Render(
			EditorQuestionTemplate,
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

func TestEditorPrompt(t *testing.T) {
	if _, err := exec.LookPath("vi"); err != nil {
		t.Skip("vi not found in PATH")
	}

	tests := []PromptTest{
		{
			"Test Editor prompt interaction",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.Send("iAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction with default",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Default: "No comment",
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.SendLine(":q!")
			},
			"No comment",
		},
		{
			"Test Editor prompt interaction overriding default",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Default: "No comment",
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.Send("iAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction hiding default",
			&Editor{
				Editor:      "vi",
				Message:     "Edit git commit message",
				Default:     "No comment",
				HideDefault: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.SendLine(":q!")
			},
			"No comment",
		},
		{
			"Test Editor prompt interaction and prompt for help",
			&Editor{
				Editor:  "vi",
				Message: "Edit git commit message",
				Help:    "Describe your git commit",
			},
			func(c *expect.Console) {
				c.ExpectString(
					fmt.Sprintf(
						"Edit git commit message [%s for help] [Enter to launch editor]",
						string(defaultPromptConfig().HelpInput),
					),
				)
				c.SendLine("?")
				c.ExpectString("Describe your git commit")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.Send("iAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
			},
			"Add editor prompt tests\n",
		},
		{
			"Test Editor prompt interaction with default and append default",
			&Editor{
				Editor:        "vi",
				Message:       "Edit git commit message",
				Default:       "No comment",
				AppendDefault: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message (No comment) [Enter to launch editor]")
				c.SendLine("")
				c.ExpectString("No comment")
				c.SendLine("dd")
				c.SendLine(":wq!")
				c.ExpectEOF()
			},
			"",
		},
		{
			"Test Editor prompt interaction with editor args",
			&Editor{
				Editor:  "vi --",
				Message: "Edit git commit message",
			},
			func(c *expect.Console) {
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				go c.ExpectEOF()
				time.Sleep(time.Millisecond)
				c.Send("iAdd editor prompt tests\x1b")
				c.SendLine(":wq!")
			},
			"Add editor prompt tests\n",
		},
	}

	for _, test := range tests {
		if s, ok := test.expected.(string); ok && s == "Add editor prompt tests\n" {
			// XXX
			test.expected = s[:len(s)-1] + "\ufeff\n"
		}
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
