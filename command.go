package app

import (
	flag "github.com/spf13/pflag"
)

type (
	// Command is a basic helper topic or when embed it can
	// be used for easily create a runnable command.
	//
	// example:
	//
	// type BarCommand struct { Command  }
	//
	// func (b BarCommand) Init(*App) error {
	// 	 return nil
	// }
	//
	// func (b BarCommand) Run([]string, *App) error {
	//	return nil
	// }
	Command struct {
		Flags *flag.FlagSet
		Usage string
		Short string
		Long  string
		Name  string
	}

	// CommandInterface are the basic methods that a command
	// should have to be a helper topic, for a runnable command
	// it should also implement CommandRunnableInterface
	CommandInterface interface {
		GetFlags() *flag.FlagSet
		GetUsage() string
		GetShort() string
		GetLong() string
		GetName() string
	}

	// CommandRunnableInterface are the basic methods to be combined
	// with the CommandInterface to be a valid runnable command
	CommandRunnableInterface interface {
		Run([]string, *App) error
		Init(*App) error
	}
)

func (c *Command) GetFlags() *flag.FlagSet { return c.Flags }
func (c *Command) GetUsage() string        { return c.Usage }
func (c *Command) GetShort() string        { return c.Short }
func (c *Command) GetLong() string         { return c.Long }
func (c *Command) GetName() string         { return c.Name }