package app

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
