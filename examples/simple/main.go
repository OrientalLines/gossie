package main

import (
	"fmt"
	"strconv"

	"github.com/orientallines/gossie/pkg/gossie"
)

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Simple Application",
		AppDescription: "This is a simple application",
		AppVersion:     "1.0.0",
		AppAuthor:      "Unknown Author",
	})

	// app.Action(func(c *gossie.Context) error {
	// 	c.Println("Hello, user!")
	// 	return nil
	// })

	app.Command("greet", func(cmd *gossie.Command) {
		cmd.Description("Greet the user")
		cmd.Action(func(c *gossie.Context) error {
			c.Println("Hello, user!")
			return nil
		})

		// define subcommand
		cmd.Command("formal", func(subcmd *gossie.Command) {
			subcmd.Description("Formal greeting")
			subcmd.Action(func(c *gossie.Context) error {
				c.Println("Good day, esteemed user.")
				return nil
			})
		})
	})

	app.Command("math", func(cmd *gossie.Command) {
		cmd.Description("Perform mathematical operations")
		cmd.Flag("verbose", "Enable verbose output").Short('v').Alias("details")

		cmd.Command("add", func(subcmd *gossie.Command) {
			subcmd.Description("Add two numbers")
			subcmd.Arg("x", "First number").Required()
			subcmd.Arg("y", "Second number").Required()
			subcmd.Action(func(c *gossie.Context) error {
				// Using named argument access
				xStr := c.GetArg("x")
				yStr := c.GetArg("y")

				if xStr == "" || yStr == "" {
					return fmt.Errorf("add command requires both x and y arguments")
				}

				x, err := strconv.ParseFloat(xStr, 64)
				if err != nil {
					return fmt.Errorf("invalid number for x: %s", xStr)
				}

				y, err := strconv.ParseFloat(yStr, 64)
				if err != nil {
					return fmt.Errorf("invalid number for y: %s", yStr)
				}

				result := x + y
				fmt.Printf("%.2f + %.2f = %.2f\n", x, y, result)
				return nil
			})
		})

		cmd.Command("subtract", func(subcmd *gossie.Command) {
			subcmd.Description("Subtract two numbers")
			subcmd.Arg("x").Required().Description("First number")
			subcmd.Arg("y").Required().Description("Second number")
			subcmd.Action(func(c *gossie.Context) error {
				// Using named argument access
				xStr := c.GetArg("x")
				yStr := c.GetArg("y")

				if xStr == "" || yStr == "" {
					return fmt.Errorf("subtract command requires both x and y arguments")
				}

				x, err := strconv.ParseFloat(xStr, 64)
				if err != nil {
					return fmt.Errorf("invalid number for x: %s", xStr)
				}

				y, err := strconv.ParseFloat(yStr, 64)
				if err != nil {
					return fmt.Errorf("invalid number for y: %s", yStr)
				}

				result := x - y
				fmt.Printf("%.2f - %.2f = %.2f\n", x, y, result)
				return nil
			})
		})

		cmd.Command("multiply", func(subcmd *gossie.Command) {
			subcmd.Description("Multiply numbers")
			subcmd.Arg("factors").Multiple().Description("Numbers to multiply")
			subcmd.Action(func(c *gossie.Context) error {
				args := c.Args()
				if len(args) < 2 {
					return fmt.Errorf("multiply command requires at least 2 arguments")
				}

				result := 1.0
				numbers := make([]float64, 0, len(args))

				for i, arg := range args {
					num, err := strconv.ParseFloat(arg, 64)
					if err != nil {
						return fmt.Errorf("invalid number: %s", arg)
					}
					numbers = append(numbers, num)
					if i == 0 {
						result = num
					} else {
						result *= num
					}
				}

				// Print the multiplication expression
				fmt.Print(numbers[0])
				for i := 1; i < len(numbers); i++ {
					fmt.Printf(" × %.2f", numbers[i])
				}
				fmt.Printf(" = %.2f\n", result)
				return nil
			})
		})
	})

	app.Command("config", func(cmd *gossie.Command) {
		cmd.Description("Manage configuration")
		cmd.Flag("verbose", "Enable verbose output").Short('v').Alias("details")

		cmd.Command("set", func(subcmd *gossie.Command) {
			subcmd.Description("Set a configuration value")
			subcmd.Arg("key").Required().Description("Configuration key")
			subcmd.Arg("value").Required().Description("Configuration value")
			subcmd.Action(func(c *gossie.Context) error {
				// Using named argument access
				key := c.GetArg("key")
				value := c.GetArg("value")

				if key == "" || value == "" {
					return fmt.Errorf("set command requires both key and value arguments")
				}

				// Check for verbose flag
				if c.HasFlag("verbose") {
					fmt.Printf("Setting configuration with verbose output: %s = %s\n", key, value)
				} else {
					fmt.Printf("Setting configuration: %s = %s\n", key, value)
				}
				return nil
			})
		})

		cmd.Command("get", func(subcmd *gossie.Command) {
			subcmd.Description("Get a configuration value")
			subcmd.Arg("key").Required().Description("Configuration key")
			subcmd.Action(func(c *gossie.Context) error {
				// Using named argument access
				key := c.GetArg("key")

				if key == "" {
					return fmt.Errorf("get command requires a key argument")
				}

				// Check for verbose flag
				if c.HasFlag("verbose") {
					fmt.Printf("Getting configuration with verbose output: %s\n", key)
				} else {
					fmt.Printf("Getting configuration: %s\n", key)
				}
				return nil
			})
		})

		cmd.Action(func(c *gossie.Context) error {
			fmt.Println("Configuration management commands:")
			fmt.Println("  config set <key> <value>  - Set a configuration value")
			fmt.Println("  config get <key>         - Get a configuration value")
			return nil
		})
	})

	app.Run()
}
