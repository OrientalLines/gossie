package gossie

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestAppInitialization tests basic app creation and configuration
func TestAppInitialization(t *testing.T) {
	config := AppConfig{
		AppName:        "TestApp",
		AppDescription: "A test application",
		AppVersion:     "1.0.0",
		AppAuthor:      "Test Author",
	}

	app := NewApp(config)

	if app.name != config.AppName {
		t.Errorf("Expected app name %s, got %s", config.AppName, app.name)
	}
	if app.description != config.AppDescription {
		t.Errorf("Expected description %s, got %s", config.AppDescription, app.description)
	}
	if app.version != config.AppVersion {
		t.Errorf("Expected version %s, got %s", config.AppVersion, app.version)
	}
	if app.author != config.AppAuthor {
		t.Errorf("Expected author %s, got %s", config.AppAuthor, app.author)
	}

	// Check that help command is auto-registered
	if _, exists := app.commands["help"]; !exists {
		t.Error("Help command should be auto-registered")
	}
}

// TestAppInitializationDefaults tests app initialization with empty config
func TestAppInitializationDefaults(t *testing.T) {
	app := NewApp(AppConfig{})

	if app.name != "Gossie CLI Application" {
		t.Errorf("Expected default name, got %s", app.name)
	}
	if app.version != "1.0.0" {
		t.Errorf("Expected default version 1.0.0, got %s", app.version)
	}
	if app.author != "Unknown Author" {
		t.Errorf("Expected default author, got %s", app.author)
	}
}

// TestCommandRegistration tests traditional command registration
func TestCommandRegistration(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	// Test command registration
	cmd := app.Command("test", func(cmd *Command) {
		cmd.Description("A test command")
	})

	if cmd == nil {
		t.Fatal("Command should not be nil")
	}

	if cmd.name != "test" {
		t.Errorf("Expected command name 'test', got '%s'", cmd.name)
	}

	if cmd.description != "A test command" {
		t.Errorf("Expected description 'A test command', got '%s'", cmd.description)
	}

	// Check command is registered in app
	if registeredCmd, exists := app.commands["test"]; !exists {
		t.Error("Command should be registered in app")
	} else if registeredCmd != cmd {
		t.Error("Registered command should be the same instance")
	}
}

// TestSubcommandRegistration tests nested command structure
func TestSubcommandRegistration(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	var executed bool
	var receivedArg string

	app.Command("parent", func(parentCmd *Command) {
		parentCmd.Description("Parent command")

		parentCmd.Command("child", func(childCmd *Command) {
			childCmd.Description("Child command")
			childCmd.Arg("name").Required().Description("Name argument")

			childCmd.Action(func(c *Context) error {
				executed = true
				receivedArg = c.GetArg("name")
				return nil
			})
		})
	})

	// Verify parent command exists
	parentCmd, exists := app.commands["parent"]
	if !exists {
		t.Fatal("Parent command should exist")
	}

	// Verify child command exists in parent
	childCmd, exists := parentCmd.subcommands["child"]
	if !exists {
		t.Fatal("Child command should exist in parent")
	}

	// Test child command execution
	childCmd.execute([]string{"test-name"})

	if !executed {
		t.Error("Child command should have been executed")
	}
	if receivedArg != "test-name" {
		t.Errorf("Expected argument 'test-name', got '%s'", receivedArg)
	}
}

// TestContextMethods tests Context functionality
func TestContextMethods(t *testing.T) {
	// Create a mock context for testing
	ctx := &Context{
		args:    []string{"arg1", "arg2", "--flag1=value1", "--flag2", "value2"},
		argMap:  map[string]string{"name": "John", "age": "25"},
		flagMap: map[string]string{"verbose": "true", "output": "json"},
	}

	// Test Args() method
	args := ctx.Args()
	if len(args) != 5 {
		t.Errorf("Expected 5 args, got %d", len(args))
	}

	// Test Arg() method
	if ctx.Arg(0) != "arg1" {
		t.Errorf("Expected arg[0] = 'arg1', got '%s'", ctx.Arg(0))
	}
	if ctx.Arg(10) != "" {
		t.Errorf("Expected arg[10] = '', got '%s'", ctx.Arg(10))
	}

	// Test GetArg() method
	if ctx.GetArg("name") != "John" {
		t.Errorf("Expected GetArg('name') = 'John', got '%s'", ctx.GetArg("name"))
	}
	if ctx.GetArg("nonexistent") != "" {
		t.Errorf("Expected GetArg('nonexistent') = '', got '%s'", ctx.GetArg("nonexistent"))
	}

	// Test GetFlag() method
	if ctx.GetFlag("verbose") != "true" {
		t.Errorf("Expected GetFlag('verbose') = 'true', got '%s'", ctx.GetFlag("verbose"))
	}
	if ctx.GetFlag("nonexistent") != "" {
		t.Errorf("Expected GetFlag('nonexistent') = '', got '%s'", ctx.GetFlag("nonexistent"))
	}

	// Test HasFlag() method
	if !ctx.HasFlag("verbose") {
		t.Error("Expected HasFlag('verbose') = true")
	}
	if ctx.HasFlag("nonexistent") {
		t.Error("Expected HasFlag('nonexistent') = false")
	}
}

// TestContextNilSafety tests that Context methods handle nil maps gracefully
func TestContextNilSafety(t *testing.T) {
	ctx := &Context{
		args: []string{"test"},
		// argMap and flagMap are nil
	}

	// These should not panic
	if ctx.GetArg("anything") != "" {
		t.Error("GetArg should return empty string for nil map")
	}
	if ctx.GetFlag("anything") != "" {
		t.Error("GetFlag should return empty string for nil map")
	}
	if ctx.HasFlag("anything") {
		t.Error("HasFlag should return false for nil map")
	}
}

// TestStructBasedCommandRegistration tests AddCommand with struct-based commands
func TestStructBasedCommandRegistration(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	// Test struct for command
	type TestCommand struct {
		Name    string `gossie:"name"`
		Verbose bool   `gossie:"verbose"`
		Count   int    `gossie:"count"`
	}

	app.AddCommand(func(cmd *TestCommand) error {
		return nil
	})

	// Check that command was registered with converted name
	expectedCmdName := "test-command" // TestCommand -> test-command
	if _, exists := app.commands[expectedCmdName]; !exists {
		t.Errorf("Command should be registered with name '%s'", expectedCmdName)
	}
}

// TestCommandExecution tests the execution flow
func TestCommandExecution(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	var executionCount int
	var lastArgs []string

	app.Command("test-exec", func(cmd *Command) {
		cmd.Description("Test execution command")
		cmd.Arg("value").Required().Description("Test value")

		cmd.Action(func(c *Context) error {
			executionCount++
			lastArgs = c.Args()
			return nil
		})
	})

	// Get the command and execute it
	cmd := app.commands["test-exec"]
	if cmd == nil {
		t.Fatal("Command should exist")
	}

	// Execute with test arguments
	cmd.execute([]string{"test-value"})

	if executionCount != 1 {
		t.Errorf("Expected execution count 1, got %d", executionCount)
	}

	if len(lastArgs) != 1 || lastArgs[0] != "test-value" {
		t.Errorf("Expected args ['test-value'], got %v", lastArgs)
	}
}

// TestCommandExecutionWithFlags tests execution with flags
func TestCommandExecutionWithFlags(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	var receivedFlags map[string]string

	app.Command("test-flags", func(cmd *Command) {
		cmd.Description("Test flags command")
		cmd.Flag("verbose", "Verbose output").Short('v')
		cmd.Flag("output", "Output format").Short('o')

		cmd.Action(func(c *Context) error {
			receivedFlags = make(map[string]string)
			for k, v := range c.flagMap {
				receivedFlags[k] = v
			}
			return nil
		})
	})

	cmd := app.commands["test-flags"]
	cmd.execute([]string{"--verbose", "--output=json"})

	if receivedFlags["verbose"] != "true" {
		t.Errorf("Expected verbose flag to be 'true', got '%s'", receivedFlags["verbose"])
	}
	if receivedFlags["output"] != "json" {
		t.Errorf("Expected output flag to be 'json', got '%s'", receivedFlags["output"])
	}
}

// TestCommandWithMultipleArguments tests commands with multiple arguments
func TestCommandWithMultipleArguments(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	var receivedArgs map[string]string

	app.Command("multi-arg", func(cmd *Command) {
		cmd.Description("Multiple arguments command")
		cmd.Arg("first").Required().Description("First argument")
		cmd.Arg("second").Required().Description("Second argument")
		cmd.Arg("third").Description("Optional third argument")

		cmd.Action(func(c *Context) error {
			receivedArgs = make(map[string]string)
			for k, v := range c.argMap {
				receivedArgs[k] = v
			}
			return nil
		})
	})

	cmd := app.commands["multi-arg"]
	cmd.execute([]string{"value1", "value2", "value3"})

	if receivedArgs["first"] != "value1" {
		t.Errorf("Expected first arg 'value1', got '%s'", receivedArgs["first"])
	}
	if receivedArgs["second"] != "value2" {
		t.Errorf("Expected second arg 'value2', got '%s'", receivedArgs["second"])
	}
	if receivedArgs["third"] != "value3" {
		t.Errorf("Expected third arg 'value3', got '%s'", receivedArgs["third"])
	}
}

// TestFlagParsing tests various flag parsing scenarios
func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "Simple boolean flags",
			args: []string{"--verbose", "--dry-run"},
			expected: map[string]string{
				"verbose": "true",
				"dry-run": "true",
			},
		},
		{
			name: "Flags with values",
			args: []string{"--output=json", "--level=debug"},
			expected: map[string]string{
				"output": "json",
				"level":  "debug",
			},
		},
		{
			name: "Mixed flags and arguments",
			args: []string{"arg1", "--flag1", "arg2", "--flag2=value"},
			expected: map[string]string{
				"flag1": "true",
				"flag2": "value",
			},
		},
		{
			name: "Short flags",
			args: []string{"-v", "-o", "json"},
			expected: map[string]string{
				"v": "true",
				"o": "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				args:    tt.args,
				argMap:  make(map[string]string),
				flagMap: make(map[string]string),
			}

			cmd := &Command{}
			cmd.parseArgsAndFlags(ctx)

			for expectedKey, expectedValue := range tt.expected {
				if actualValue, exists := ctx.flagMap[expectedKey]; !exists {
					t.Errorf("Expected flag '%s' to exist", expectedKey)
				} else if actualValue != expectedValue {
					t.Errorf("Expected flag '%s' = '%s', got '%s'", expectedKey, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestPascalToKebabConversion tests struct name to command name conversion
func TestPascalToKebabConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"TestCommand", "test-command"},
		{"HTTPClient", "h-t-t-p-client"},
		{"APIHandler", "a-p-i-handler"},
		{"SimpleTest", "simple-test"},
		{"VeryLongCommandName", "very-long-command-name"},
		{"X", "x"},
		{"ABC", "a-b-c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			app := NewApp(AppConfig{})
			// Create a struct type with the test name
			var structType reflect.Type
			switch tt.input {
			case "TestCommand":
				type TestCommand struct{}
				structType = reflect.TypeOf(TestCommand{})
			case "HTTPClient":
				type HTTPClient struct{}
				structType = reflect.TypeOf(HTTPClient{})
			case "APIHandler":
				type APIHandler struct{}
				structType = reflect.TypeOf(APIHandler{})
			case "SimpleTest":
				type SimpleTest struct{}
				structType = reflect.TypeOf(SimpleTest{})
			case "VeryLongCommandName":
				type VeryLongCommandName struct{}
				structType = reflect.TypeOf(VeryLongCommandName{})
			case "X":
				type X struct{}
				structType = reflect.TypeOf(X{})
			case "ABC":
				type ABC struct{}
				structType = reflect.TypeOf(ABC{})
			}
			result := app.getCommandNameFromStruct(structType)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestGossieTagParsing tests the tag parsing functionality
func TestGossieTagParsing(t *testing.T) {
	tests := []struct {
		tag      string
		expected struct {
			name string
			opts map[string]string
		}
	}{
		{
			tag: "name",
			expected: struct {
				name string
				opts map[string]string
			}{name: "name", opts: map[string]string{}},
		},
		{
			tag: "verbose,short=v",
			expected: struct {
				name string
				opts map[string]string
			}{name: "verbose", opts: map[string]string{"short": "v"}},
		},
		{
			tag: "count,required",
			expected: struct {
				name string
				opts map[string]string
			}{name: "count", opts: map[string]string{"required": "true"}},
		},
		{
			tag: "output,short=o,alias=out",
			expected: struct {
				name string
				opts map[string]string
			}{name: "output", opts: map[string]string{"short": "o", "alias": "out"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			app := NewApp(AppConfig{})
			name, opts := app.parseGossieTag(tt.tag)

			if name != tt.expected.name {
				t.Errorf("Expected name '%s', got '%s'", tt.expected.name, name)
			}

			for expectedKey, expectedValue := range tt.expected.opts {
				if actualValue, exists := opts[expectedKey]; !exists {
					t.Errorf("Expected option '%s' to exist", expectedKey)
				} else if actualValue != expectedValue {
					t.Errorf("Expected option '%s' = '%s', got '%s'", expectedKey, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestErrorHandling tests error handling in commands
func TestErrorHandling(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	app.Command("error-test", func(cmd *Command) {
		cmd.Description("Command that always returns an error")
		cmd.Action(func(c *Context) error {
			return fmt.Errorf("this is a test error")
		})
	})

	// Capture stdout to check error output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := app.commands["error-test"]
	cmd.execute([]string{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "Error: this is a test error") {
		t.Errorf("Expected error message in output, got: %s", output)
	}
}

// TestRequiredArguments tests validation of required arguments
func TestRequiredArguments(t *testing.T) {
	app := NewApp(AppConfig{AppName: "TestApp"})

	var executionCount int

	app.Command("required-test", func(cmd *Command) {
		cmd.Description("Command with required argument")
		cmd.Arg("name").Required().Description("Required name")

		cmd.Action(func(c *Context) error {
			// Check if required argument is provided
			if c.GetArg("name") == "" {
				return fmt.Errorf("name argument is required")
			}
			executionCount++
			return nil
		})
	})

	cmd := app.commands["required-test"]

	// Test with missing argument
	cmd.execute([]string{})
	if executionCount != 0 {
		t.Error("Command should not execute with missing required argument")
	}

	// Test with argument provided
	cmd.execute([]string{"test-name"})
	if executionCount != 1 {
		t.Error("Command should execute with required argument provided")
	}
}
