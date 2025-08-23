package gossie

import (
	"fmt"
	"reflect"
	"testing"
)

// TestParseFlagFromArgs tests flag parsing from command line arguments
func TestParseFlagFromArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flagName string
		expected string
	}{
		{
			name:     "Simple boolean flag",
			args:     []string{"--verbose"},
			flagName: "verbose",
			expected: "true",
		},
		{
			name:     "Flag with equals value",
			args:     []string{"--output=json"},
			flagName: "output",
			expected: "json",
		},
		{
			name:     "Flag with space value",
			args:     []string{"--output", "json"},
			flagName: "output",
			expected: "json",
		},
		{
			name:     "Short flag with space value",
			args:     []string{"-o", "json"},
			flagName: "o",
			expected: "json",
		},
		{
			name:     "Flag not present",
			args:     []string{"--verbose"},
			flagName: "output",
			expected: "",
		},
		{
			name:     "Flag with multiple equals",
			args:     []string{"--config=key=value"},
			flagName: "config",
			expected: "key=value",
		},
		{
			name:     "Flag with empty value",
			args:     []string{"--output="},
			flagName: "output",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(AppConfig{})
			result := app.parseFlagFromArgs(tt.args, tt.flagName)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestGetArgFields tests extraction of argument fields from struct
func TestGetArgFields(t *testing.T) {
	type TestStruct struct {
		Name    string `gossie:"name"`
		Verbose bool   `gossie:"verbose"`
		Count   int    `gossie:"count"`
		Hidden  string `not-a-gossie-tag`
		NoTag   string
	}

	app := NewApp(AppConfig{})
	structType := reflect.TypeOf(TestStruct{})
	argFields := app.getArgFields(structType)

	// Should find Name and Count as argument fields (not Verbose which is a flag)
	expectedNames := map[string]bool{
		"Name":  true,
		"Count": true,
	}

	if len(argFields) != 2 {
		t.Errorf("Expected 2 argument fields, got %d", len(argFields))
	}

	for _, field := range argFields {
		if !expectedNames[field.Name] {
			t.Errorf("Unexpected argument field: %s", field.Name)
		}
	}
}

// TestIsFieldFlag tests determination of flag vs argument fields
func TestIsFieldFlag(t *testing.T) {
	tests := []struct {
		name      string
		fieldType reflect.Type
		tag       string
		expected  bool
	}{
		{
			name:      "Bool field should be flag",
			fieldType: reflect.TypeOf(true),
			tag:       "verbose",
			expected:  true,
		},
		{
			name:      "String field should be argument",
			fieldType: reflect.TypeOf(""),
			tag:       "name",
			expected:  false,
		},
		{
			name:      "Int field should be argument",
			fieldType: reflect.TypeOf(0),
			tag:       "count",
			expected:  false,
		},
		{
			name:      "Explicit flag tag",
			fieldType: reflect.TypeOf(""),
			tag:       "output,flag=true",
			expected:  true,
		},
		{
			name:      "Explicit argument tag",
			fieldType: reflect.TypeOf(true),
			tag:       "dry-run,arg=true",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(AppConfig{})
			result := app.isFieldFlag(tt.fieldType, tt.tag)
			if result != tt.expected {
				t.Errorf("Expected isFieldFlag to be %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestTypeConversion tests conversion of string values to different types
func TestTypeConversion(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldType reflect.Type
		expected  interface{}
		shouldErr bool
	}{
		{
			name:      "String to string",
			value:     "hello",
			fieldType: reflect.TypeOf(""),
			expected:  "hello",
			shouldErr: false,
		},
		{
			name:      "String to int valid",
			value:     "42",
			fieldType: reflect.TypeOf(0),
			expected:  int64(42),
			shouldErr: false,
		},
		{
			name:      "String to int invalid",
			value:     "not-a-number",
			fieldType: reflect.TypeOf(0),
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "String to bool true",
			value:     "true",
			fieldType: reflect.TypeOf(true),
			expected:  true,
			shouldErr: false,
		},
		{
			name:      "String to bool false",
			value:     "false",
			fieldType: reflect.TypeOf(true),
			expected:  false,
			shouldErr: false,
		},
		{
			name:      "String to bool invalid",
			value:     "maybe",
			fieldType: reflect.TypeOf(true),
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "String to float valid",
			value:     "3.14",
			fieldType: reflect.TypeOf(0.0),
			expected:  3.14,
			shouldErr: false,
		},
		{
			name:      "String to float invalid",
			value:     "not-a-float",
			fieldType: reflect.TypeOf(0.0),
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(AppConfig{})

			// Create a field value to test
			fieldValue := reflect.New(tt.fieldType).Elem()

			err := app.setFieldValue(fieldValue, tt.value)

			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			var actual interface{}
			switch tt.fieldType.Kind() {
			case reflect.String:
				actual = fieldValue.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				actual = fieldValue.Int()
			case reflect.Bool:
				actual = fieldValue.Bool()
			case reflect.Float32, reflect.Float64:
				actual = fieldValue.Float()
			}

			if actual != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

// TestComplexFlagParsing tests complex flag parsing scenarios
func TestComplexFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected map[string]string
		leftover []string
	}{
		{
			name:     "Flags before arguments",
			args:     []string{"--verbose", "--output=json", "arg1", "arg2"},
			expected: map[string]string{"verbose": "true", "output": "json"},
			leftover: []string{"arg1", "arg2"},
		},
		{
			name:     "Flags after arguments",
			args:     []string{"arg1", "--verbose", "arg2"},
			expected: map[string]string{"verbose": "true"},
			leftover: []string{"arg1", "arg2"},
		},
		{
			name: "Mixed flags and arguments",
			args: []string{"--flag1", "arg1", "--flag2=value", "arg2", "--flag3"},
			expected: map[string]string{
				"flag1": "true",
				"flag2": "value",
				"flag3": "true",
			},
			leftover: []string{"arg1", "arg2"},
		},
		{
			name:     "Only flags",
			args:     []string{"--verbose", "--dry-run", "--output=console"},
			expected: map[string]string{"verbose": "true", "dry-run": "true", "output": "console"},
			leftover: []string{},
		},
		{
			name:     "Only arguments",
			args:     []string{"arg1", "arg2", "arg3"},
			expected: map[string]string{},
			leftover: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				args:    make([]string, len(tt.args)),
				argMap:  make(map[string]string),
				flagMap: make(map[string]string),
			}
			copy(ctx.args, tt.args)

			cmd := &Command{}
			cmd.parseArgsAndFlags(ctx)

			// Check flags
			for expectedKey, expectedValue := range tt.expected {
				if actualValue, exists := ctx.flagMap[expectedKey]; !exists {
					t.Errorf("Expected flag '%s' to exist", expectedKey)
				} else if actualValue != expectedValue {
					t.Errorf("Expected flag '%s' = '%s', got '%s'", expectedKey, expectedValue, actualValue)
				}
			}

			// Check leftover arguments
			if len(ctx.args) != len(tt.leftover) {
				t.Errorf("Expected %d leftover args, got %d", len(tt.leftover), len(ctx.args))
			}

			for i, expectedArg := range tt.leftover {
				if i >= len(ctx.args) || ctx.args[i] != expectedArg {
					t.Errorf("Expected leftover arg[%d] = '%s', got '%s'", i, expectedArg, ctx.args[i])
				}
			}
		})
	}
}

// TestStructFieldParsing tests parsing struct fields into command configuration
func TestStructFieldParsing(t *testing.T) {
	type TestCommand struct {
		Name     string   `gossie:"name,required"`
		Verbose  bool     `gossie:"verbose,short=v"`
		Count    int      `gossie:"count"`
		Output   string   `gossie:"output,flag=true"`
		Multiple []string `gossie:"files,multiple"`
	}

	app := NewApp(AppConfig{})
	structType := reflect.TypeOf(TestCommand{})

	// Create a mock command to test field addition
	cmd := &Command{
		args:  make([]*Argument, 0),
		flags: make([]*Flag, 0),
	}

	// Add fields to command
	app.addFieldsToCommandTest(cmd, structType)

	// Verify arguments were added correctly
	expectedArgs := map[string]*Argument{}
	for _, arg := range cmd.args {
		expectedArgs[arg.name] = arg
	}

	if len(expectedArgs) != 3 {
		t.Errorf("Expected 3 arguments, got %d", len(expectedArgs))
	}

	// Check name argument (required)
	if nameArg, exists := expectedArgs["name"]; !exists {
		t.Error("Expected 'name' argument to exist")
	} else if !nameArg.required {
		t.Error("Expected 'name' argument to be required")
	}

	// Check count argument
	if _, exists := expectedArgs["count"]; !exists {
		t.Error("Expected 'count' argument to exist")
	}

	// Check files argument (multiple)
	if filesArg, exists := expectedArgs["files"]; !exists {
		t.Error("Expected 'files' argument to exist")
	} else if !filesArg.multiple {
		t.Error("Expected 'files' argument to allow multiple values")
	}

	// Verify flags were added correctly
	expectedFlags := map[string]*Flag{}
	for _, flag := range cmd.flags {
		expectedFlags[flag.name] = flag
	}

	if len(expectedFlags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(expectedFlags))
	}

	// Check verbose flag (with short name)
	if verboseFlag, exists := expectedFlags["verbose"]; !exists {
		t.Error("Expected 'verbose' flag to exist")
	} else if verboseFlag.shortName != 'v' {
		t.Errorf("Expected 'verbose' flag short name 'v', got '%c'", verboseFlag.shortName)
	}

	// Check output flag (explicitly marked as flag)
	if _, exists := expectedFlags["output"]; !exists {
		t.Error("Expected 'output' flag to exist")
	}
}

// Helper method to add fields to command for testing
func (app *App) addFieldsToCommandTest(cmd *Command, structType reflect.Type) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("gossie")
		if tag == "" {
			continue
		}

		fieldName, options := app.parseGossieTag(tag)

		if app.isFieldFlag(field.Type, tag) {
			// Add as flag
			flag := &Flag{
				name:        fieldName,
				description: fmt.Sprintf("Set %s", fieldName),
			}
			if shortName := options["short"]; shortName != "" && len(shortName) == 1 {
				flag.shortName = rune(shortName[0])
			}
			cmd.flags = append(cmd.flags, flag)
		} else {
			// Add as argument
			arg := &Argument{
				name:        fieldName,
				description: fmt.Sprintf("Set %s", fieldName),
			}
			if options["required"] == "true" {
				arg.required = true
			}
			if options["multiple"] == "true" {
				arg.multiple = true
			}
			cmd.args = append(cmd.args, arg)
		}
	}
}
