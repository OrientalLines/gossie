package main

import (
	"fmt"

	"github.com/orientallines/gossie/pkg/gossie"
)

// Example 1: Basic Command Registration
func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Basic CLI",
		AppDescription: "Demonstrates basic Gossie framework usage",
		AppVersion:     "1.0.0",
		AppAuthor:      "Gossie Examples",
	})

	// Example 1: Simple command with traditional approach
	app.Command("hello", func(cmd *gossie.Command) {
		cmd.Description("Prints a greeting message")
		cmd.Action(func(c *gossie.Context) error {
			fmt.Println("Hello, World!")
			return nil
		})
	})

	// Example 2: Command with arguments using traditional approach
	app.Command("greet", func(cmd *gossie.Command) {
		cmd.Description("Greet a person by name")
		cmd.Arg("name").Required().Description("Name to greet")
		cmd.Action(func(c *gossie.Context) error {
			name := c.GetArg("name")
			if name == "" {
				return fmt.Errorf("name argument is required")
			}
			fmt.Printf("Hello, %s! 👋\n", name)
			return nil
		})
	})

	// Example 3: Command with flags using traditional approach
	app.Command("info", func(cmd *gossie.Command) {
		cmd.Description("Display information with optional formatting")
		cmd.Flag("verbose", "Show detailed information").Short('v')
		cmd.Flag("uppercase", "Convert output to uppercase").Short('u')
		cmd.Action(func(c *gossie.Context) error {
			message := "This is a basic CLI application built with Gossie framework"

			if c.HasFlag("verbose") {
				message += "\nFramework: Gossie CLI v1.0.0"
				message += "\nAuthor: Gossie Team"
				message += "\nPurpose: Demonstration of basic CLI patterns"
			}

			if c.HasFlag("uppercase") {
				// Convert to uppercase (simple implementation)
				bytes := []byte(message)
				for i, b := range bytes {
					if b >= 'a' && b <= 'z' {
						bytes[i] = b - 32
					}
				}
				message = string(bytes)
			}

			fmt.Println(message)
			return nil
		})
	})

	// Example 4: Struct-based command registration
	type EchoCommand struct {
		Message string `gossie:"message"`
		Repeat  int    `gossie:"repeat"`
		Reverse bool   `gossie:"reverse"`
	}

	app.AddCommand(func(cmd *EchoCommand) error {
		if cmd.Message == "" {
			return fmt.Errorf("message is required")
		}

		// Default repeat to 1 if not specified
		if cmd.Repeat <= 0 {
			cmd.Repeat = 1
		}

		// Process the message
		message := cmd.Message
		if cmd.Reverse {
			runes := []rune(message)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			message = string(runes)
		}

		// Repeat the message
		for i := 0; i < cmd.Repeat; i++ {
			fmt.Printf("Echo %d: %s\n", i+1, message)
		}

		return nil
	})

	app.Run()
}
