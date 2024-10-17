package main

import (
	"fmt"

	"github.com/orientallines/gossie/pkg/gossie"
)

type GreetCommand struct {
	Name string `gossie:"name"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Tags Application",
		AppDescription: "This is a simple application",
		AppVersion:     "1.0.0",
		AppAuthor:      "Unknown Author",
	})

	app.AddCommand(func(c *GreetCommand) error {
		fmt.Printf("Hello, %s!\n", c.Name)
		return nil
	})
}
