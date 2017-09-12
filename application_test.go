package app

import (
	"os"
	"fmt"

	flag "github.com/spf13/pflag"
)

func ExampleApplication_help_failed() {
	app := NewApp(
		newRunnableTestCommand("foo"),
		newRunnableTestCommand("bar"),
	)
	app.bin = "test_app"
	if err := app.help([]string{"help", "foobar"}...); err != nil {
		fmt.Println(err.Error())
	}
	//Output:
	//
	// usage: test_app help command
	//
	// Too many arguments given.

}

func ExampleApplication_help() {
	app := NewApp(
		newRunnableTestCommand("foo"),
		newRunnableTestCommand("bar"),
	)
	app.bin = "test_app"
	if err := app.Run([]string{"help", "foo"}); err != nil {
		fmt.Println(err.Error())
	}
	//Output:
	//
	// usage: test_app foo [-c config] [--force] [--dry-run] [ARGS]...
	//
	// The is a simple noop command for testing purpose.
}

func ExampleApplication_failed() {
	app := NewApp(
		newRunnableTestCommand("foo"),
		newRunnableTestCommand("bar"),
	)
	app.bin = "test_app"
	if err := app.Run([]string{"not:exist:cmd"}); err != nil {
		fmt.Println(err.Error())
	}
	//Output:
	//
	// Unknown subcommand "not:exist:cmd".
	// Run 'test_app help' for usage.
}

func ExampleApplication_usage() {
	app := NewApp(
		newRunnableTestCommand("foo"),
		newRunnableTestCommand("bar"),
		newTestCommand("help:bar"),
		newTestCommand("help:foo"),
	)
	app.bin = "test_app"
	app.Usage(os.Stdout)
	//Output:
	//
	// Usage:
	// 	test_app command [arguments]
	//
	// The commands are:
	//
	// 	foo         This a one line info for foo
	// 	bar         This a one line info for bar
	//
	// Use "test_app help [command]" for more information about a command.
	//
	// Additional help topics:
	//
	// 	help:bar    This a one line info for help:bar
	// 	help:foo    This a one line info for help:foo
	//
	// Use "test_app help [topic]" for more information about that topic.
}

func ExampleApplication_usage_extra() {
	app := NewApp(
		newRunnableTestCommand("foo"),
		newRunnableTestCommand("bar"),
		newTestCommand("help:bar"),
		newTestCommand("help:foo"),
	)
	app.Intro = `Some example intro.........`
	app.bin = "test_app"
	app.Usage(os.Stdout)
	//Output:
	// Some example intro.........
	// Usage:
	// 	test_app command [arguments]
	//
	// The commands are:
	//
	// 	foo         This a one line info for foo
	// 	bar         This a one line info for bar
	//
	// Use "test_app help [command]" for more information about a command.
	//
	// Additional help topics:
	//
	// 	help:bar    This a one line info for help:bar
	// 	help:foo    This a one line info for help:foo
	//
	// Use "test_app help [topic]" for more information about that topic.
}

func ExampleApplication_usage_simple() {
	app := NewApp(newRunnableTestCommand("foo"), newRunnableTestCommand("bar"))
	app.bin = "test_app"
	app.Usage(os.Stdout)
	//Output:
	//
	// Usage:
	// 	test_app command [arguments]
	//
	// The commands are:
	//
	// 	foo         This a one line info for foo
	// 	bar         This a one line info for bar
	//
	// Use "test_app help [command]" for more information about a command.
}

func ExampleApplication_usage_basic() {
	app := NewApp()
	app.bin = "test_app"
	app.Usage(os.Stdout)
	//Output:
	//
	// Usage:
	// 	test_app command [arguments]
	//
}

type noOpCommand struct { CommandInterface }
func (t noOpCommand) Run([]string, *App) error { return nil }
func (t noOpCommand) Init(*App) error { return nil }

func newRunnableTestCommand(n string) CommandInterface {
	return &noOpCommand{newTestCommand(n)}
}

func newTestCommand(n string) CommandInterface {
	return &Command{
		Flags: new(flag.FlagSet),
		Usage: "[-c config] [--force] [--dry-run] [ARGS]...",
		Short: "This a one line info for " + n,
		Long:  "The is a simple noop command for testing purpose.",
		Name:  n,
	}
}