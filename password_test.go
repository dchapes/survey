package survey

import (
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	expect "github.com/Netflix/go-expect"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestPasswordRender(t *testing.T) {
	tests := []struct {
		title    string
		prompt   Password
		data     PasswordTemplateData
		expected string
	}{
		{
			"Test Password question output",
			Password{Message: "Tell me your secret:"},
			PasswordTemplateData{},
			fmt.Sprintf("%s Tell me your secret: ", defaultIcons().Question.Text),
		},
		{
			"Test Password question output with help hidden",
			Password{Message: "Tell me your secret:", Help: "This is helpful"},
			PasswordTemplateData{},
			fmt.Sprintf("%s Tell me your secret: [%s for help] ", defaultIcons().Question.Text, defaultPromptConfig().HelpInput),
		},
		{
			"Test Password question output with help shown",
			Password{Message: "Tell me your secret:", Help: "This is helpful"},
			PasswordTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s Tell me your secret: ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
	}

	for _, test := range tests {
		test.data.Password = test.prompt

		// set the icon set
		test.data.Config = defaultPromptConfig()

		actual, err := core.RunTemplate(
			PasswordQuestionTemplate,
			&test.data,
		)
		if err != nil {
			t.Errorf("%s failed:\n\t%v", test.title, err)
			continue
		}
		if actual != test.expected {
			t.Errorf("%s:\n\tgave %q\n\twant %q", test.title, actual, test.expected)
		}
	}
}

func TestPasswordPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Password prompt interaction",
			&Password{
				Message: "Please type your password",
			},
			func(c *expect.Console) {
				c.ExpectString("Please type your password")
				c.Send("secret")
				c.SendLine("")
				c.ExpectEOF()
			},
			"secret",
		},
		{
			"Test Password prompt interaction with help",
			&Password{
				Message: "Please type your password",
				Help:    "It's a secret",
			},
			func(c *expect.Console) {
				c.ExpectString("Please type your password")
				c.SendLine("?")
				c.ExpectString("It's a secret")
				c.Send("secret")
				c.SendLine("")
				c.ExpectEOF()
			},
			"secret",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
