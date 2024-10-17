package gossie

import (
	"fmt"
)

type Command struct {
	name        string
	description string
	usage       string
	author      string
	version     string
	shortName   rune
	help        string
	action      func(*Context) error
	subcommands map[string]*Command
	args        []*Argument
	flags       []*Flag
	app         *App
}

func (c *Command) Description(desc string) {
	c.description = desc
}

func (c *Command) Action(fn func(*Context) error) {
	c.action = fn
}

func (c *Command) Command(name string, configFn func(*Command)) *Command {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command)
	}
	cmd := &Command{
		name: name,
		app:  c.app,
	}
	configFn(cmd)
	c.subcommands[name] = cmd
	return cmd
}

func (c *Command) Arg(name string, desc ...string) *Argument {
	arg := &Argument{
		name: name,
	}
	if len(desc) > 0 {
		arg.description = desc[0]
	}
	c.args = append(c.args, arg)
	return arg
}

func (c *Command) Flag(name, desc string) *Flag {
	flag := &Flag{
		name:        name,
		description: desc,
	}
	c.flags = append(c.flags, flag)
	return flag
}

func (c *Command) execute(args []string) {
	// Check for help flag
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		c.printHelp()
		return
	}

	// Check for subcommands
	if len(args) > 0 {
		if subCmd, ok := c.subcommands[args[0]]; ok {
			subCmd.execute(args[1:])
			return
		}
	}

	// Execute the command's action
	ctx := &Context{
		command: c,
		args:    args,
	}
	if c.action != nil {
		err := c.action(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			c.printUsage()
		}
	} else {
		fmt.Printf("Command '%s' has no action defined\n", c.name)
		c.printUsage()
	}
}

// Add this new method
func (c *Command) printHelp() {
	printCommandHelp(c)
}

// Add a new method to print usage
func (c *Command) printUsage() {
	printUsage(c.app.name+" "+c.name, c.usage)
}

type Argument struct {
	name        string
	description string
	required    bool
	multiple    bool
}

func (a *Argument) Required() *Argument {
	a.required = true
	return a
}

func (a *Argument) Multiple() *Argument {
	a.multiple = true
	return a
}

func (a *Argument) Description(desc string) *Argument {
	a.description = desc
	return a
}

type Flag struct {
	name        string
	shortName   rune
	description string
	aliases     []string
	isBool      bool
}

func (f *Flag) Short(short rune) *Flag {
	f.shortName = short
	return f
}

func (f *Flag) Alias(alias string) *Flag {
	f.aliases = append(f.aliases, alias)
	return f
}

type Context struct {
	command *Command
	args    []string
}

func (c *Context) Println(a ...interface{}) {
	fmt.Println(a...)
}
