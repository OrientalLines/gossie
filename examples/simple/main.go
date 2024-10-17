package main

import "github.com/orientallines/gossie/pkg/gossie"

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

		cmd.Command("add", func(subcmd *gossie.Command) {
			subcmd.Description("Add two numbers")
			subcmd.Arg("x", "First number").Required()
			subcmd.Arg("y", "Second number").Required()
			subcmd.Action(func(c *gossie.Context) error {
				// Implement addition logic here
				return nil
			})
		})

		cmd.Command("subtract", func(subcmd *gossie.Command) {
			subcmd.Description("Subtract two numbers")
			subcmd.Arg("x").Required().Description("First number")
			subcmd.Arg("y").Required().Description("Second number")
			subcmd.Action(func(c *gossie.Context) error {
				// Implement subtraction logic here
				return nil
			})
		})

		cmd.Command("multiply", func(subcmd *gossie.Command) {
			subcmd.Description("Multiply numbers")
			subcmd.Arg("factors").Multiple().Description("Numbers to multiply")
			subcmd.Action(func(c *gossie.Context) error {
				// Implement multiplication logic here
				return nil
			})
		})
	})

	app.Command("config", func(cmd *gossie.Command) {
		cmd.Description("Manage configuration")
		cmd.Flag("verbose", "Enable verbose output").Short('v').Alias("details")
		cmd.Action(func(c *gossie.Context) error {
			// Implement config logic here
			return nil
		})
	})

	app.Run()
}
