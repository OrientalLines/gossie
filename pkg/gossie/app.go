package gossie

import (
	"fmt"
	"reflect"
)

type App struct {
	appName        string
	appDescription string
	appVersion     string
	appAuthor      string
}

type AppConfig struct {
	AppName        string
	AppDescription string
	AppVersion     string
	AppAuthor      string
}

// NewApp creates a new Gossie CLI Application
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

	return &App{
		appName:        config.AppName,
		appDescription: config.AppDescription,
		appVersion:     config.AppVersion,
		appAuthor:      config.AppAuthor,
	}

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

func (app *App) Parse() {
	return
}
