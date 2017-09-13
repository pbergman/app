package app

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

type App struct {
	Container 	interface{}
	temperating *template.Template
	commands    []CommandInterface
	Name        string
	Intro       string
	bin         string
	StdOut      io.Writer
	StdErr      io.Writer
}

// Usage will print the application usage to the given writer
func (a *App) Usage(w io.Writer) {
	a.tmpl(UsageTmpl, w, a.commands)
}

func (a *App) tmpl(t string, w io.Writer, data interface{}) {
	template.Must(template.Must(a.GetTemplateEngine().Clone()).Parse(t)).Execute(w, data)
}

// GetCommand will search and return the command by the given name
func (a *App) GetCommand(name string) CommandInterface {
	for _, cmd := range a.commands {
		if cmd.GetName() == name {
			return cmd
		}
	}
	return nil
}

// Run will run the command with the given args, if a error is returned it
// will be of a app.Error and has a exit function.
//
// app := NewApp(....)
// if err := app.Run(os.Args); err != nil {
//    err.(Error).Exit(os.Stderr)
// }
func (a *App) Run(args []string) error {
	switch name := args[0]; name {
	case "help":
		return a.help(args[1:]...)
	default:
		if cmd := a.GetCommand(name); cmd != nil {
			if runnable, ok := cmd.(CommandRunnableInterface); ok {
				if err := runnable.Init(a); err != nil {
					return Error{err, 3}
				}
				flags := cmd.GetFlags()
				setFlagsUsage(flags, func() { a.help(cmd.GetName()); os.Exit(2) })
				flags.Parse(args[1:])
				if err := runnable.Run(cmd.GetFlags().Args(), a); err != nil {
					return Error{err, 4}
				} else {
					return nil
				}
			}
		}
	}
	return Error{fmt.Errorf("Unknown subcommand %q.\nRun '%s help' for usage.\n", args[0], a.bin), 2}
}

func (a *App) help(args ...string) error {
	if len(args) == 0 {
		a.Usage(a.StdOut)
		return nil
	}
	if len(args) != 1 {
		return &Error{fmt.Errorf("usage: %s help command\n\nToo many arguments given.\n", a.bin), 5}
	}
	if cmd := a.GetCommand(args[0]); cmd != nil {
		if _, o := cmd.(CommandRunnableInterface); o {
			fmt.Fprintf(a.StdOut, "usage: %s %s %s\n\n", a.bin, cmd.GetName(), cmd.GetUsage())
		}
		a.tmpl(cmd.GetLong(), a.StdOut, nil)
	}
	return nil
}

// GetTemplateEngine will return the template engine and add the default functions
func (a *App) GetTemplateEngine() *template.Template {
	if a.temperating == nil {
		a.temperating = template.New("app").Funcs(
			template.FuncMap{
				"runnable": func(c interface{}) bool {
					_, o := c.(CommandRunnableInterface)
					return o
				},
				"has_runnable": func(list []CommandInterface) bool {
					for _, c := range list {
						if _, ok := c.(CommandRunnableInterface); ok {
							return true
						}
					}
					return false
				},
				"has_not_runnable": func(list []CommandInterface) bool {
					for _, c := range list {
						if _, ok := c.(CommandRunnableInterface); !ok {
							return true
						}
					}
					return false
				},
				"exec_bin": func() string {
					return a.bin
				},
				"has_intro": func() bool {
					return len(a.Intro) > 0
				},
				"intro": func() string {
					return a.Intro
				},
			},
		)
	}
	return a.temperating
}

func NewApp(command ...CommandInterface) *App {
	return &App{
		commands: command,
		bin:      os.Args[0],
		StdOut:   os.Stdout,
		StdErr:   os.Stderr,
	}
}
