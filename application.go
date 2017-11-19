package app

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"strings"
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
	PreRun 		func(CommandInterface) error
}

// Usage will print the application usage to the given writer
func (a *App) Usage(w io.Writer) {
	a.tmpl(UsageTmpl, w, a.getCommandList())
}

func (a *App) getCommandList() *commandList {
	list := &commandList{0, make(map[string][]CommandInterface), make(map[string][]CommandInterface)}
	for i, c := 0, len(a.commands); i < c; i++ {
		name  := a.commands[i].GetName()
		group := DEFAULT_GROUP_NAME
		if s := len(name); s > list.Max {
			list.Max = s
		}
		if index := strings.Index(name, ":"); index > 0 {
			group = name[:index]
		}
		if _, ok := a.commands[i].(CommandRunnableInterface); ok {
			if _, ok := list.Runnable[group]; !ok {
				list.Runnable[group] = make([]CommandInterface, 0)
			}
			list.Runnable[group] = append(list.Runnable[group], a.commands[i])
		}  else {
			if _, ok := list.Helpers[group]; !ok {
				list.Helpers[group] = make([]CommandInterface, 0)
			}
			list.Helpers[group] = append(list.Helpers[group], a.commands[i])
		}
	}
	return list
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
			if a.PreRun != nil {
				if err := a.PreRun(cmd); err != nil {
					return &Error{err, 6}
				}
			}
			if runnable, ok := cmd.(CommandRunnableInterface); ok {
				if err := runnable.Init(a); err != nil {
					return &Error{err, 3}
				}
				if flags := cmd.GetFlags(); flags != nil {
					setFlagsUsage(flags, func() { a.help(cmd.GetName()); os.Exit(2) })
					if err := flags.Parse(args[1:]); err != nil {
						return &Error{err, 5}
					}
					if err := runnable.Run(cmd.GetFlags().Args(), a); err != nil {
						return &Error{err, 4}
					} else {
						return nil
					}
				}
			}
		}
	}
	return &Error{fmt.Errorf("Unknown subcommand %q.\nRun '%s help' for usage.\n", args[0], a.bin), 2}
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
		a.tmpl(cmd.GetLong(), a.StdOut, a.Container)
	}
	return nil
}

// GetTemplateEngine will return the template engine and add the default functions
func (a *App) GetTemplateEngine() *template.Template {
	if a.temperating == nil {
		a.temperating = template.New("app").Funcs(
			template.FuncMap{
				"space": func(name string, max int) string {
					size := (max - len(name)) + 2
					buf := make([]byte, size)
					for i := 0; i < size; i++ {
						buf[i] = ' '
					}
					return string(buf)
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
