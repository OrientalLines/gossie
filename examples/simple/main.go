package main

import "github.com/orientallines/gossie/pkg/gossie"

type Args struct {
	name  string `short:"n" default:"SomeName"`
	count uint8  `short:"c" default:"5"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Simple Application",
		AppDescription: "This is a simple application",
		AppVersion:     "1.0.0",
		AppAuthor:      "Unknown Author",
	})

	app.AddArgs(Args{})
	app.Parse()
}
