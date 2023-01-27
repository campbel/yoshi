package yoshi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = Create[TestYoshiApp]("appy")

func TestHelp(t *testing.T) {
	testHelp(t, ctx, "", `Usage: appy [options] COMMAND
Commands:
  call
  email
  text
Options:
  -n string (default: smudge)
  -a string
`)
}

func TestHelpSubCommandCall(t *testing.T) {
	testHelp(t, ctx, "call", `Usage: appy call [options] COMMAND
Commands:
  message
Options:
  -n int The number to call
`)
}

func TestHelpSubCommandCallMessage(t *testing.T) {
	testHelp(t, ctx, "call message", `Usage: appy call message [options]
Options:
  -t string The text to send
`)
}

func TestHelpSubCommandEmail(t *testing.T) {
	testHelp(t, ctx, "email", `Usage: appy email [options]
Options:
  -a string
`)
}

func TestHelpSubCommandText(t *testing.T) {
	testHelp(t, ctx, "text", `Usage: appy text [options]
Options:
  -n string
`)
}

func TestHelpErrors(t *testing.T) {
	t.Run("unknown command", func(t *testing.T) {
		testHelp(t, ctx, "trousers", "Usage: appy [options] COMMAND\nCommands:\n  call\n  email\n  text\nOptions:\n  -n string (default: smudge)\n  -a string\n")
	})
	t.Run("unknown subcommand", func(t *testing.T) {
		testHelp(t, ctx, "call trousers", "Usage: appy call [options] COMMAND\nCommands:\n  message\nOptions:\n  -n int The number to call\n")
	})
}

func testHelp[T any](t *testing.T, ctx *Context[T], args, out string) {
	t.Helper()
	assert := assert.New(t)
	assert.Equal(out, ctx.help(strings.Split(args, " ")...))
}
