package survey

import (
	"fmt"
	"reflect"
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

func Stdio(c *expect.Console) terminal.Stdio {
	return terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
}

type PromptTest struct {
	name      string
	prompt    Prompt
	procedure func(*expect.Console)
	expected  interface{}
}

func RunPromptTest(t *testing.T, test PromptTest) {
	t.Helper()
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		var err error
		if p, ok := test.prompt.(wantsStdio); ok {
			p.WithStdio(stdio)
		}

		answer, err = test.prompt.Prompt(defaultPromptConfig())
		return err
	})
	if !reflect.DeepEqual(answer, test.expected) {
		t.Fatalf("%s\n\tgave: %#v\n\twant: %#v", test.name, answer, test.expected)
	}
}

func TestPagination(t *testing.T) {
	tests := []struct {
		choices   []string
		pageSize  int
		sel       int
		i, j, idx int
	}{
		{ // too few options
			choices:  []string{"choice1", "choice2", "choice3"},
			pageSize: 4,
			sel:      3,
			i:        0, j: 3, idx: 3,
		},
		{ // first half
			choices:  []string{"choice1", "choice2", "choice3", "choice4", "choice5", "choice6"},
			pageSize: 4,
			sel:      2,
			i:        0, j: 4, idx: 2,
		},
		{ // middle
			choices:  []string{"choice0", "choice1", "choice2", "choice3", "choice4", "choice5"},
			pageSize: 2,
			sel:      3,
			i:        2, j: 4, idx: 1,
		},
		{ // last half
			choices:  []string{"choice0", "choice1", "choice2", "choice3", "choice4", "choice5"},
			pageSize: 3,
			sel:      5,
			i:        3, j: 6, idx: 2,
		},
	}
	for _, tc := range tests {
		choices := core.OptionAnswerList(tc.choices)
		page, idx := paginate(tc.pageSize, choices, tc.sel)
		want := choices[tc.i:tc.j]
		if idx != tc.idx || !reflect.DeepEqual(page, want) {
			t.Errorf("paginate(%d, %v, %d)\n\tgave: %v, %d\n\twant: %v, %d",
				tc.pageSize, choices, tc.sel,
				page, idx,
				choices[tc.i:tc.j], tc.idx,
			)
			continue
		}
	}
}

func TestAsk(t *testing.T) {
	t.Skip()
	tests := []struct {
		name      string
		questions []*Question
		procedure func(*expect.Console)
		expected  map[string]interface{}
	}{
		{
			"Test Ask for all prompts",
			[]*Question{
				{
					Name: "pizza",
					Prompt: &Confirm{
						Message: "Is pizza your favorite food?",
					},
				},
				{
					Name: "commit-message",
					Prompt: &Editor{
						Message: "Edit git commit message",
					},
				},
				{
					Name: "commit-message-validated",
					Prompt: &Editor{
						Message: "Edit git commit message",
					},
					Validate: func(v interface{}) error {
						s := v.(string)
						if strings.Contains(s, "invalid") {
							return fmt.Errorf("invalid error message")
						}
						return nil
					},
				},
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
				},
				{
					Name: "day",
					Prompt: &MultiSelect{
						Message: "What days do you prefer:",
						Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
					},
				},
				{
					Name: "password",
					Prompt: &Password{
						Message: "Please type your password",
					},
				},
				{
					Name: "color",
					Prompt: &Select{
						Message: "Choose a color:",
						Options: []string{"red", "blue", "green", "yellow"},
					},
				},
			},
			func(c *expect.Console) {
				// Confirm
				c.ExpectString("Is pizza your favorite food? (y/N)")
				c.SendLine("Y")

				// Editor
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("iAdd editor prompt tests\x1b")
				c.SendLine(":wq!")

				// Editor validated
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.Send("i invalid input first try\x1b")
				c.SendLine(":wq!")
				time.Sleep(time.Millisecond)
				c.ExpectString("invalid error message")
				c.ExpectString("Edit git commit message [Enter to launch editor]")
				c.SendLine("")
				time.Sleep(time.Millisecond)
				c.ExpectString("first try")
				c.Send("ccAdd editor prompt tests\x1b")
				c.SendLine(":wq!")

				// Input
				c.ExpectString("What is your name?")
				c.SendLine("Johnny Appleseed")

				// MultiSelect
				c.ExpectString("What days do you prefer:  [Use arrows to move, space to select, type to filter]")
				// Select Monday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(" ")
				// Select Wednesday.
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine(" ")

				// Password
				c.ExpectString("Please type your password")
				c.Send("secret")
				c.SendLine("")

				// Select
				c.ExpectString("Choose a color:  [Use arrows to move, type to filter]")
				c.SendLine("yellow")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"pizza":                    true,
				"commit-message":           "Add editor prompt tests\n",
				"commit-message-validated": "Add editor prompt tests\n",
				"name":                     "Johnny Appleseed",
				"day":                      []string{"Monday", "Wednesday"},
				"password":                 "secret",
				"color":                    "yellow",
			},
		},
		{
			"Test Ask with validate survey.Required",
			[]*Question{
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
					Validate: Required,
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectString("Sorry, your reply was invalid: Value is required")
				c.SendLine("Johnny Appleseed")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"name": "Johnny Appleseed",
			},
		},
		{
			"Test Ask with transformer survey.ToLower",
			[]*Question{
				{
					Name: "name",
					Prompt: &Input{
						Message: "What is your name?",
					},
					Transform: ToLower,
				},
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Johnny Appleseed")
				c.ExpectEOF()
			},
			map[string]interface{}{
				"name": "johnny appleseed",
			},
		},
	}

	for _, test := range tests {
		// Capture range variable.
		test := test
		t.Run(test.name, func(t *testing.T) {
			answers := make(map[string]interface{})
			RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
				return Ask(test.questions, &answers, WithStdio(stdio.In, stdio.Out, stdio.Err))
			})
			if !reflect.DeepEqual(answers, test.expected) {
				t.Fatalf("%s\n\tgave: %v\n\twant: %v", test.name, answers, test.expected)
			}
		})
	}
}

func TestAsk_returnsErrorIfTargetIsNil(t *testing.T) {
	// pass an empty place to leave the answers
	err := Ask([]*Question{}, nil)

	// if we didn't get an error
	if err == nil {
		// the test failed
		t.Error("Did not encounter error when asking with no where to record.")
	}
}
