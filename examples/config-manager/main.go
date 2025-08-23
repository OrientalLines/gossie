package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/orientallines/gossie/pkg/gossie"
)

type ConfigSetCommand struct {
	Key     string `gossie:"key"`
	Value   string `gossie:"value"`
	Type    string `gossie:"type"`
	Section string `gossie:"section"`
	Global  bool   `gossie:"global"`
}

type ConfigGetCommand struct {
	Key     string `gossie:"key"`
	Section string `gossie:"section"`
	All     bool   `gossie:"all"`
}

type ConfigListCommand struct {
	Section string `gossie:"section"`
	Global  bool   `gossie:"global"`
}

type ConfigUnsetCommand struct {
	Key     string `gossie:"key"`
	Section string `gossie:"section"`
	All     bool   `gossie:"all"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Config Manager",
		AppDescription: "A powerful configuration management CLI",
		AppVersion:     "1.0.0",
		AppAuthor:      "Config Team",
	})

	// Set configuration value using struct-based approach
	app.AddCommand(func(cmd *ConfigSetCommand) error {
		if cmd.Key == "" || cmd.Value == "" {
			return fmt.Errorf("both key and value are required")
		}

		// Validate type if specified
		if cmd.Type != "" {
			validTypes := []string{"string", "int", "bool", "float"}
			isValid := false
			for _, t := range validTypes {
				if cmd.Type == t {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid type: %s. Valid types: %s", cmd.Type, strings.Join(validTypes, ", "))
			}
		}

		// Validate value based on type
		if cmd.Type != "" {
			switch cmd.Type {
			case "int":
				if _, err := strconv.Atoi(cmd.Value); err != nil {
					return fmt.Errorf("value must be a valid integer: %s", cmd.Value)
				}
			case "bool":
				if _, err := strconv.ParseBool(cmd.Value); err != nil {
					return fmt.Errorf("value must be a valid boolean: %s", cmd.Value)
				}
			case "float":
				if _, err := strconv.ParseFloat(cmd.Value, 64); err != nil {
					return fmt.Errorf("value must be a valid float: %s", cmd.Value)
				}
			}
		}

		section := "default"
		if cmd.Section != "" {
			section = cmd.Section
		}

		scope := "local"
		if cmd.Global {
			scope = "global"
		}

		fmt.Printf("Setting configuration [%s] %s.%s = %s (%s)\n",
			scope, section, cmd.Key, cmd.Value, cmd.Type)

		// Simulate saving to file
		fmt.Printf("✅ Configuration saved to ~/.config/app/config.json\n")
		return nil
	})

	// Get configuration value using struct-based approach
	app.AddCommand(func(cmd *ConfigGetCommand) error {
		if cmd.All {
			fmt.Println("Current configuration:")
			fmt.Println("=====================")

			if cmd.Section != "" {
				fmt.Printf("[%s]\n", cmd.Section)
			} else {
				fmt.Println("[default]")
			}

			fmt.Println("  app.name = MyApp")
			fmt.Println("  app.version = 1.0.0")
			fmt.Println("  database.host = localhost")
			fmt.Println("  database.port = 5432")
			fmt.Println("  database.ssl = true")
			fmt.Println("  logging.level = info")
			fmt.Println("  logging.file = /var/log/app.log")
		} else if cmd.Key != "" {

			// Mock configuration values
			mockValues := map[string]string{
				"app.name":      "MyApp",
				"app.version":   "1.0.0",
				"database.host": "localhost",
				"database.port": "5432",
				"database.ssl":  "true",
				"logging.level": "info",
				"logging.file":  "/var/log/app.log",
				"server.port":   "8080",
				"server.host":   "0.0.0.0",
				"cache.enabled": "true",
				"cache.ttl":     "3600",
			}

			fullKey := cmd.Key
			if cmd.Section != "" {
				fullKey = cmd.Section + "." + cmd.Key
			}

			if value, exists := mockValues[fullKey]; exists {
				fmt.Printf("%s=%s\n", fullKey, value)
			} else {
				fmt.Printf("Configuration key '%s' not found\n", fullKey)
			}
		} else {
			return fmt.Errorf("either --key or --all must be specified")
		}

		return nil
	})

	// List configuration using struct-based approach
	app.AddCommand(func(cmd *ConfigListCommand) error {
		fmt.Println("Configuration files:")
		fmt.Println("====================")

		if cmd.Global {
			fmt.Println("Global configuration (~/.config/app/config.json):")
		} else {
			fmt.Println("Local configuration (.config.json):")
		}

		config := map[string]interface{}{
			"app": map[string]interface{}{
				"name":    "MyApp",
				"version": "1.0.0",
				"debug":   false,
			},
			"database": map[string]interface{}{
				"host": "localhost",
				"port": 5432,
				"ssl":  true,
				"pool": map[string]interface{}{
					"min": 5,
					"max": 20,
				},
			},
			"server": map[string]interface{}{
				"host": "0.0.0.0",
				"port": 8080,
				"cors": map[string]interface{}{
					"enabled": true,
					"origins": []string{"*"},
				},
			},
			"logging": map[string]interface{}{
				"level":     "info",
				"file":      "/var/log/app.log",
				"max_size":  100,
				"max_files": 5,
			},
		}

		if cmd.Section != "" {
			if sectionConfig, exists := config[cmd.Section]; exists {
				jsonData, _ := json.MarshalIndent(sectionConfig, "", "  ")
				fmt.Println(string(jsonData))
			} else {
				return fmt.Errorf("section '%s' not found", cmd.Section)
			}
		} else {
			jsonData, _ := json.MarshalIndent(config, "", "  ")
			fmt.Println(string(jsonData))
		}

		return nil
	})

	// Unset configuration using struct-based approach
	app.AddCommand(func(cmd *ConfigUnsetCommand) error {
		if cmd.All {
			fmt.Println("⚠️  This will remove ALL configuration values!")
			fmt.Println("Are you sure? (This action cannot be undone)")
			return fmt.Errorf("use --force to confirm")
		}

		if cmd.Key == "" {
			return fmt.Errorf("--key is required when not using --all")
		}

		section := "default"
		if cmd.Section != "" {
			section = cmd.Section
		}

		fmt.Printf("Removing configuration: %s.%s\n", section, cmd.Key)
		fmt.Println("✅ Configuration removed successfully")

		return nil
	})

	// Traditional commands for more complex operations
	app.Command("validate", func(cmd *gossie.Command) {
		cmd.Description("Validate configuration files")
		cmd.Arg("files").Multiple().Description("Configuration files to validate")
		cmd.Flag("strict", "Enable strict validation").Short('s')
		cmd.Flag("verbose", "Show detailed validation results").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			files := c.Args()

			if len(files) == 0 {
				// Check default config file
				fmt.Println("Validating default configuration...")
				fmt.Println("✅ ~/.config/app/config.json - Valid")
				return nil
			}

			for _, file := range files {
				fmt.Printf("Validating %s...\n", file)

				if _, err := os.Stat(file); os.IsNotExist(err) {
					return fmt.Errorf("configuration file not found: %s", file)
				}

				// Mock validation
				if c.HasFlag("strict") {
					fmt.Printf("✅ %s - Valid (strict mode)\n", file)
				} else {
					fmt.Printf("✅ %s - Valid\n", file)
				}

				if c.HasFlag("verbose") {
					fmt.Printf("  - Syntax: OK\n")
					fmt.Printf("  - Schema: OK\n")
					fmt.Printf("  - References: OK\n")
				}
			}

			return nil
		})
	})

	app.Command("diff", func(cmd *gossie.Command) {
		cmd.Description("Compare configuration files")
		cmd.Arg("file1").Required().Description("First configuration file")
		cmd.Arg("file2").Required().Description("Second configuration file")
		cmd.Flag("unified", "Show unified diff format").Short('u')
		cmd.Flag("context", "Show context diff format").Short('c')
		cmd.Action(func(c *gossie.Context) error {
			file1 := c.GetArg("file1")
			file2 := c.GetArg("file2")

			if c.HasFlag("unified") {
				fmt.Printf("--- %s\n", file1)
				fmt.Printf("+++ %s\n", file2)
				fmt.Println("@@ -1,4 +1,4 @@")
				fmt.Println(" database.host = localhost")
				fmt.Println("-database.port = 5432")
				fmt.Println("+database.port = 3306")
				fmt.Println(" logging.level = info")
				fmt.Println("+logging.file = /tmp/app.log")
			} else if c.HasFlag("context") {
				fmt.Printf("*** %s\n", file1)
				fmt.Printf("--- %s\n", file2)
				fmt.Println("***************")
				fmt.Println("*** 1,4 ****")
				fmt.Println("  database.host = localhost")
				fmt.Println("  database.port = 5432")
				fmt.Println("  logging.level = info")
				fmt.Println("--- 1,5 ----")
				fmt.Println("  database.host = localhost")
				fmt.Println("  database.port = 3306")
				fmt.Println("  logging.level = info")
				fmt.Println("  logging.file = /tmp/app.log")
			} else {
				fmt.Println("Differences between configuration files:")
				fmt.Println("=====================================")
				fmt.Printf("- %s: database.port = 5432\n", file1)
				fmt.Printf("+ %s: database.port = 3306\n", file2)
				fmt.Printf("+ %s: logging.file = /tmp/app.log\n", file2)
			}

			return nil
		})
	})

	app.Command("backup", func(cmd *gossie.Command) {
		cmd.Description("Backup configuration files")
		cmd.Arg("destination").Description("Backup destination directory")
		cmd.Flag("compress", "Compress the backup").Short('c')
		cmd.Flag("include-global", "Include global configuration")
		cmd.Action(func(c *gossie.Context) error {
			destination := c.GetArg("destination")
			if destination == "" {
				destination = "./config-backup"
			}

			fmt.Printf("Creating configuration backup...\n")

			if c.HasFlag("include-global") {
				fmt.Println("Including global configuration files...")
			}

			if c.HasFlag("compress") {
				fmt.Printf("Creating compressed backup: %s.tar.gz\n", destination)
			} else {
				fmt.Printf("Creating backup directory: %s\n", destination)
			}

			fmt.Println("Backing up local configuration...")
			fmt.Println("Backing up environment-specific configs...")

			if c.HasFlag("include-global") {
				fmt.Println("Backing up global configuration...")
			}

			fmt.Printf("✅ Configuration backup completed: %s\n", destination)
			return nil
		})
	})

	app.Run()
}
