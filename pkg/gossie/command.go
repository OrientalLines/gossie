package gossie

import (
	"fmt"
	"strings"
)

type Command struct {
	name        string
	description string
	usage       string
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

	// Create context with parsed arguments and flags
	ctx := &Context{
		command: c,
		args:    args,
		argMap:  make(map[string]string),
		flagMap: make(map[string]string),
	}

	// Parse arguments and flags into maps
	c.parseArgsAndFlags(ctx)

	// Execute the command's action
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
	argMap  map[string]string // parsed arguments by name
	flagMap map[string]string // parsed flags by name
}

func (c *Context) Println(a ...interface{}) {
	fmt.Println(a...)
}

// Args returns the command line arguments
func (c *Context) Args() []string {
	return c.args
}

// Arg returns the argument at the specified index, or empty string if index is out of bounds
func (c *Context) Arg(index int) string {
	if index < 0 || index >= len(c.args) {
		return ""
	}
	return c.args[index]
}

// ArgMap returns the parsed arguments map
func (c *Context) ArgMap() map[string]string {
	return c.argMap
}

// FlagMap returns the parsed flags map
func (c *Context) FlagMap() map[string]string {
	return c.flagMap
}

// GetArg returns the argument value by name, or empty string if not found
func (c *Context) GetArg(name string) string {
	if c.argMap == nil {
		return ""
	}
	return c.argMap[name]
}

// GetFlag returns the flag value by name, or empty string if not found
func (c *Context) GetFlag(name string) string {
	if c.flagMap == nil {
		return ""
	}
	return c.flagMap[name]
}

// HasFlag returns true if the flag is present
func (c *Context) HasFlag(name string) bool {
	if c.flagMap == nil {
		return false
	}
	_, exists := c.flagMap[name]
	return exists
}

// parseArgsAndFlags parses command line arguments and flags into the context maps
func (c *Command) parseArgsAndFlags(ctx *Context) {
	args := ctx.args
	argIndex := 0

	// First pass: parse flags
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check if it's a flag
		if strings.HasPrefix(arg, "--") {
			// Long flag
			flagName := strings.TrimPrefix(arg, "--")
			if strings.Contains(flagName, "=") {
				// --flag=value format
				parts := strings.SplitN(flagName, "=", 2)
				flagName = parts[0]
				flagValue := parts[1]
				ctx.flagMap[flagName] = flagValue
			} else {
				// Check if this flag expects a value by looking ahead
				if i+1 < len(args) {
					nextArg := args[i+1]
					// If next arg doesn't start with dash, it might be a value
					// But only if it's not obviously a flag name
					if !strings.HasPrefix(nextArg, "-") {
						// Check if the next argument looks like a flag name (contains common flag patterns)
						if !strings.Contains(nextArg, ".") && !strings.Contains(nextArg, "/") && !strings.Contains(nextArg, ":") {
							// This could be a flag value, but it's ambiguous
							// For now, treat boolean flags as boolean
							ctx.flagMap[flagName] = "true"
						} else {
							// This looks like a positional argument (contains ., :, /)
							ctx.flagMap[flagName] = "true"
						}
					} else {
						// Next arg starts with dash, so this is a boolean flag
						ctx.flagMap[flagName] = "true"
					}
				} else {
					// No more arguments, this is a boolean flag
					ctx.flagMap[flagName] = "true"
				}
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Short flag
			flagName := arg[1:]
			if len(flagName) == 1 {
				// Single character flag
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					// -f value format
					ctx.flagMap[flagName] = args[i+1]
					i++ // Skip the next argument as it's the flag value
				} else {
					// -f format (boolean flag)
					ctx.flagMap[flagName] = "true"
				}
			} else {
				// Multiple character short flag (treat as boolean for now)
				for _, char := range flagName {
					ctx.flagMap[string(char)] = "true"
				}
			}
		} else {
			// This is a positional argument
			ctx.args[argIndex] = arg
			argIndex++
		}
	}

	// Trim the args slice to remove processed flags
	ctx.args = ctx.args[:argIndex]

	// Second pass: map positional arguments to argument names
	argPos := 0
	for _, argument := range c.args {
		if argPos < len(ctx.args) {
			ctx.argMap[argument.name] = ctx.args[argPos]
			argPos++
		}
	}
}
