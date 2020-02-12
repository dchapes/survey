// +build !windows

package survey

import (
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
)

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Parallel()

	// Multiplex output to a buffer as well for the raw bytes.
	var sb strings.Builder
	c, state, err := vt10x.NewVT10XConsole(expect.WithStdout(&sb))
	if err != nil {
		t.Fatal("vt10x.NewVT10XConsole", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			t.Error("vt10x Close():", err)
		}
	}()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(c)
	}()

	if err := test(Stdio(c)); err != nil {
		t.Fatal(err)
	}

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	t.Logf("Raw output: %q", sb.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
