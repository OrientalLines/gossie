package gossie

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type App struct {
	application
}

type AppConfig struct {
	AppName        string
	AppDescription string
	AppVersion     string
	AppAuthor      string
}

// NewApp creates a new Gossie CLI Application with an auto-generated help command
func NewApp(config AppConfig) *App {
	if config.AppName == "" {
		config.AppName = "Gossie CLI Application"
	}
	if config.AppVersion == "" {
		config.AppVersion = "1.0.0"
	}
	if config.AppAuthor == "" {
		config.AppAuthor = "Unknown Author"
	}

	app := &App{
		application: application{
			name:          config.AppName,
			description:   config.AppDescription,
			version:       config.AppVersion,
			author:        config.AppAuthor,
			commands:      make(map[string]*Command),
			defaultAction: nil,
		},
	}

	app.addDefaultHelpCommand()

	return app
}

// addDefaultHelpCommand adds a default help command to the application
func (app *App) addDefaultHelpCommand() {
	// Check if a help command already exists
	if _, exists := app.commands["help"]; exists {
		return
	}

	app.Command("help", func(cmd *Command) {
		cmd.Description("Display help information")
		cmd.Action(func(c *Context) error {
			app.printHelp()
			return nil
		})
	})
}

// printHelp prints the help information for the application
func (app *App) printHelp() {
	printAppHelp(app)
}

// Command method remains unchanged
func (app *App) Command(name string, configFn func(*Command)) *Command {
	cmd := &Command{
		name: name,
		app:  app,
	}
	configFn(cmd)
	app.commands[cmd.name] = cmd
	return cmd
}

func (app *App) AddArgs(args ...interface{}) {
	for _, arg := range args {
		t := reflect.TypeOf(arg)

		if t.Kind() != reflect.Struct {
			fmt.Printf("Skipping non-struct argument: %v\n", arg)
			continue
		}

		fmt.Printf("Struct: %s\n", t.Name())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("  Field: %s, Type: %s\n", field.Name, field.Type)

			if len(field.Tag) > 0 {
				fmt.Printf("    Tags: %s\n", field.Tag)
			}
		}
		fmt.Println()
	}
}

func (app *App) Action(action func(*Context) error) {
	app.defaultAction = action
}

func (app *App) executeDefaultAction() {
	if app.defaultAction == nil {
		fmt.Printf("No default action defined for %s.\n\nUse '%s --help' for available commands and options.\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		return
	}
}

func (app *App) Run() {
	if len(os.Args) == 1 {
		// If no arguments are provided, execute the default action
		app.executeDefaultAction()
		return
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		app.printHelp()
		return
	}

	if os.Args[1] == "--version" || os.Args[1] == "-V" {
		fmt.Printf("%s %s\n", app.name, app.version)
		return
	}

	cmdName := os.Args[1]
	cmd, ok := app.commands[cmdName]
	if !ok {
		fmt.Printf("Unknown command: %s\n", cmdName)
		app.printHelp()
		return
	}

	cmd.execute(os.Args[2:])
}
