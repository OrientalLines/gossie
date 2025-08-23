package gossie

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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

// AddCommand adds a struct-based command to the application
func (app *App) AddCommand(handler interface{}) {
	app.addStructCommand(handler)
}

// addStructCommand processes a function that takes a struct pointer and returns an error
func (app *App) addStructCommand(handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		fmt.Printf("AddCommand: expected function, got %s\n", handlerType.Kind())
		return
	}

	if handlerType.NumIn() != 1 || handlerType.NumOut() != 1 {
		fmt.Printf("AddCommand: handler function must have exactly 1 input and 1 output parameter\n")
		return
	}

	inputType := handlerType.In(0)
	outputType := handlerType.Out(0)

	if inputType.Kind() != reflect.Ptr || inputType.Elem().Kind() != reflect.Struct {
		fmt.Printf("AddCommand: handler function must take a pointer to a struct\n")
		return
	}

	if outputType != reflect.TypeOf((*error)(nil)).Elem() {
		fmt.Printf("AddCommand: handler function must return an error\n")
		return
	}

	structType := inputType.Elem()
	commandName := app.getCommandNameFromStruct(structType)

	app.Command(commandName, func(cmd *Command) {
		cmd.Description(fmt.Sprintf("Execute %s command", commandName))
		app.addFieldsToCommand(cmd, structType)

		cmd.Action(func(c *Context) error {
			// Create a new instance of the struct
			structValue := reflect.New(structType)

			// Parse arguments and flags into the struct
			if err := app.parseArgsIntoStruct(c, structValue); err != nil {
				return err
			}

			// Call the handler function
			handlerValue := reflect.ValueOf(handler)
			results := handlerValue.Call([]reflect.Value{structValue})
			if !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil
		})
	})
}

// getCommandNameFromStruct extracts command name from struct name (convert PascalCase to kebab-case)
func (app *App) getCommandNameFromStruct(structType reflect.Type) string {
	name := structType.Name()
	// Convert PascalCase to kebab-case
	var result []rune
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '-')
		}
		if r >= 'A' && r <= 'Z' {
			result = append(result, r-'A'+'a')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// addFieldsToCommand analyzes struct fields and adds them as arguments or flags
func (app *App) addFieldsToCommand(cmd *Command, structType reflect.Type) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Check for gossie tag
		tag := field.Tag.Get("gossie")
		if tag == "" {
			continue
		}

		// Parse tag to get field name and options
		fieldName, options := app.parseGossieTag(tag)

		if app.isFieldFlag(field.Type, tag) {
			// Add as flag
			flag := cmd.Flag(fieldName, fmt.Sprintf("Set %s", fieldName))
			if shortName := options["short"]; shortName != "" && len(shortName) == 1 {
				flag.Short(rune(shortName[0]))
			}
			if aliases := options["alias"]; aliases != "" {
				for _, alias := range strings.Split(aliases, ",") {
					flag.Alias(strings.TrimSpace(alias))
				}
			}
		} else {
			// Add as argument
			arg := cmd.Arg(fieldName, fmt.Sprintf("Set %s", fieldName))
			if options["required"] == "true" {
				arg.Required()
			}
			if options["multiple"] == "true" {
				arg.Multiple()
			}
		}
	}
}

// parseGossieTag parses gossie tag into field name and options
func (app *App) parseGossieTag(tag string) (string, map[string]string) {
	parts := strings.Split(tag, ",")
	fieldName := strings.TrimSpace(parts[0])
	options := make(map[string]string)

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				options[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		} else {
			options[part] = "true"
		}
	}

	return fieldName, options
}

// isFieldFlag determines if a field should be treated as a flag based on its type and tag
func (app *App) isFieldFlag(fieldType reflect.Type, tag string) bool {
	// Parse the tag to check for explicit type hints
	_, options := app.parseGossieTag(tag)

	// If explicitly marked as argument, it's not a flag
	if options["argument"] == "true" || options["arg"] == "true" {
		return false
	}

	// If explicitly marked as flag, it is a flag
	if options["flag"] == "true" {
		return true
	}

	// Default behavior: bool fields are flags, others are arguments
	switch fieldType.Kind() {
	case reflect.Bool:
		return true
	default:
		return false
	}
}

// parseArgsIntoStruct parses CLI arguments and flags into the struct
func (app *App) parseArgsIntoStruct(c *Context, structValue reflect.Value) error {
	structType := structValue.Elem().Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Elem().Field(i)

		tag := field.Tag.Get("gossie")
		if tag == "" {
			continue
		}

		fieldName, _ := app.parseGossieTag(tag)

		if app.isFieldFlag(field.Type, tag) {
			// Handle as flag
			if err := app.setFieldFromFlag(c, fieldName, fieldValue); err != nil {
				return err
			}
		} else {
			// Handle as argument - pass the struct type for context
			if err := app.setFieldFromArgWithContext(c, fieldName, fieldValue, structType); err != nil {
				return err
			}
		}
	}

	return nil
}

// setFieldFromFlag sets a struct field from a CLI flag
func (app *App) setFieldFromFlag(c *Context, fieldName string, fieldValue reflect.Value) error {
	// Parse flags from context args
	flagValue := app.parseFlagFromArgs(c.args, fieldName)
	if flagValue == "" {
		return nil // Flag not provided, use default value
	}

	return app.setFieldValue(fieldValue, flagValue)
}

// setFieldFromArgWithContext sets a struct field from a CLI argument with struct context
func (app *App) setFieldFromArgWithContext(c *Context, fieldName string, fieldValue reflect.Value, structType reflect.Type) error {
	// For struct-based commands, we use positional arguments
	// Map field names to positional indices
	fieldIndex := -1

	// Find the index of this field among all gossie-tagged fields
	argFields := app.getArgFields(structType)
	for i, field := range argFields {
		tag := field.Tag.Get("gossie")
		if tag != "" {
			fname, _ := app.parseGossieTag(tag)
			if fname == fieldName {
				fieldIndex = i
				break
			}
		}
	}

	if fieldIndex == -1 || fieldIndex >= len(c.args) {
		return fmt.Errorf("missing required argument: %s", fieldName)
	}

	argValue := c.args[fieldIndex]
	return app.setFieldValue(fieldValue, argValue)
}

// parseFlagFromArgs parses a flag value from command line arguments
func (app *App) parseFlagFromArgs(args []string, flagName string) string {
	for i, arg := range args {
		// Check for --flag=value format
		if strings.HasPrefix(arg, "--"+flagName+"=") {
			return strings.TrimPrefix(arg, "--"+flagName+"=")
		}

		// Check for --flag value format
		if arg == "--"+flagName {
			if i+1 < len(args) {
				return args[i+1]
			}
			// Boolean flag without value
			return "true"
		}

		// Check for -f value format (assuming single character flags)
		if len(flagName) == 1 && arg == "-"+flagName {
			if i+1 < len(args) {
				return args[i+1]
			}
			// Boolean flag without value
			return "true"
		}
	}
	return ""
}

// getArgFields returns all fields that should be treated as arguments (non-flag fields)
func (app *App) getArgFields(structType reflect.Type) []reflect.StructField {
	var argFields []reflect.StructField

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if tag := field.Tag.Get("gossie"); tag != "" {
			if !app.isFieldFlag(field.Type, tag) {
				argFields = append(argFields, field)
			}
		}
	}

	return argFields
}

// setFieldValue sets a field value with type conversion
func (app *App) setFieldValue(fieldValue reflect.Value, strValue string) error {
	if !fieldValue.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(strValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := app.parseInt(strValue); err == nil {
			fieldValue.SetInt(val)
		} else {
			return fmt.Errorf("invalid integer value: %s", strValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := app.parseUint(strValue); err == nil {
			fieldValue.SetUint(val)
		} else {
			return fmt.Errorf("invalid unsigned integer value: %s", strValue)
		}
	case reflect.Float32, reflect.Float64:
		if val, err := app.parseFloat(strValue); err == nil {
			fieldValue.SetFloat(val)
		} else {
			return fmt.Errorf("invalid float value: %s", strValue)
		}
	case reflect.Bool:
		if val, err := app.parseBool(strValue); err == nil {
			fieldValue.SetBool(val)
		} else {
			return fmt.Errorf("invalid boolean value: %s", strValue)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}

	return nil
}

// parseInt parses string to int64
func (app *App) parseInt(s string) (int64, error) {
	return parseInt(s)
}

// parseUint parses string to uint64
func (app *App) parseUint(s string) (uint64, error) {
	return parseUint(s)
}

// parseFloat parses string to float64
func (app *App) parseFloat(s string) (float64, error) {
	return parseFloat(s)
}

// parseBool parses string to bool
func (app *App) parseBool(s string) (bool, error) {
	return parseBool(s)
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

// parseInt parses string to int64
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// parseUint parses string to uint64
func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// parseFloat parses string to float64
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// parseBool parses string to bool
func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
