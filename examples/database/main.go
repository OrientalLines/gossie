package main

import (
	"encoding/json"
	"fmt"

	"github.com/orientallines/gossie/pkg/gossie"
)

// DatabaseConnection represents a database connection configuration
type DatabaseConnection struct {
	Host     string `gossie:"host"`
	Port     int    `gossie:"port"`
	Username string `gossie:"username"`
	Password string `gossie:"password"`
	Database string `gossie:"database"`
	SSL      bool   `gossie:"ssl"`
}

// QueryCommand represents a database query operation
type QueryCommand struct {
	Connection DatabaseConnection
	Query      string `gossie:"query"`
	Table      string `gossie:"table"`
	Limit      int    `gossie:"limit"`
	Format     string `gossie:"format"`
	Verbose    bool   `gossie:"verbose"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "Database CLI",
		AppDescription: "A powerful database management CLI tool",
		AppVersion:     "1.0.0",
		AppAuthor:      "Database Team",
	})

	// Connect command using struct-based approach
	app.AddCommand(func(cmd *DatabaseConnection) error {
		fmt.Printf("Connecting to database...\n")
		fmt.Printf("Host: %s\n", cmd.Host)
		fmt.Printf("Port: %d\n", cmd.Port)
		fmt.Printf("Database: %s\n", cmd.Database)
		fmt.Printf("Username: %s\n", cmd.Username)
		fmt.Printf("SSL: %t\n", cmd.SSL)

		// Simulate connection
		fmt.Printf("✅ Successfully connected to %s@%s:%d/%s\n",
			cmd.Username, cmd.Host, cmd.Port, cmd.Database)

		return nil
	})

	// Query command using struct-based approach
	app.AddCommand(func(cmd *QueryCommand) error {
		if cmd.Query == "" && cmd.Table == "" {
			return fmt.Errorf("either --query or --table must be specified")
		}

		fmt.Printf("Executing query on database...\n")

		if cmd.Verbose {
			fmt.Printf("Connection: %s@%s:%d/%s\n",
				cmd.Connection.Username, cmd.Connection.Host,
				cmd.Connection.Port, cmd.Connection.Database)
		}

		var query string
		if cmd.Query != "" {
			query = cmd.Query
		} else {
			query = fmt.Sprintf("SELECT * FROM %s", cmd.Table)
			if cmd.Limit > 0 {
				query += fmt.Sprintf(" LIMIT %d", cmd.Limit)
			}
		}

		if cmd.Verbose {
			fmt.Printf("Query: %s\n", query)
		}

		// Simulate query execution
		fmt.Printf("Query executed successfully\n")

		// Generate mock results
		if cmd.Format == "json" {
			results := map[string]interface{}{
				"status": "success",
				"rows": []map[string]interface{}{
					{"id": 1, "name": "John Doe", "email": "john@example.com"},
					{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
				},
				"count": 2,
			}
			jsonData, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("+----+------------+-------------------+")
			fmt.Println("| ID | Name       | Email             |")
			fmt.Println("+----+------------+-------------------+")
			fmt.Println("| 1  | John Doe   | john@example.com  |")
			fmt.Println("| 2  | Jane Smith | jane@example.com  |")
			fmt.Println("+----+------------+-------------------+")
			fmt.Printf("2 rows returned\n")
		}

		return nil
	})

	// Traditional command approach for more complex operations
	app.Command("schema", func(cmd *gossie.Command) {
		cmd.Description("Database schema operations")

		cmd.Command("show", func(subcmd *gossie.Command) {
			subcmd.Description("Show database schema")
			subcmd.Arg("table").Description("Specific table to show")
			subcmd.Flag("all", "Show all tables").Short('a')
			subcmd.Flag("verbose", "Show detailed schema information").Short('v')
			subcmd.Action(func(c *gossie.Context) error {
				table := c.GetArg("table")

				if c.HasFlag("all") || table == "" {
					fmt.Println("Tables in database:")
					fmt.Println("  users")
					fmt.Println("  products")
					fmt.Println("  orders")
					fmt.Println("  categories")
				} else {
					if c.HasFlag("verbose") {
						fmt.Printf("Schema for table '%s':\n", table)
						fmt.Println("CREATE TABLE users (")
						fmt.Println("  id INTEGER PRIMARY KEY,")
						fmt.Println("  name VARCHAR(255) NOT NULL,")
						fmt.Println("  email VARCHAR(255) UNIQUE,")
						fmt.Println("  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
						fmt.Println(");")
					} else {
						fmt.Printf("Columns in table '%s':\n", table)
						fmt.Println("  id (INTEGER)")
						fmt.Println("  name (VARCHAR)")
						fmt.Println("  email (VARCHAR)")
						fmt.Println("  created_at (TIMESTAMP)")
					}
				}

				return nil
			})
		})

		cmd.Command("create", func(subcmd *gossie.Command) {
			subcmd.Description("Create a new table")
			subcmd.Arg("name").Required().Description("Table name")
			subcmd.Arg("columns").Multiple().Description("Column definitions")
			subcmd.Flag("if-not-exists", "Only create if table doesn't exist")
			subcmd.Action(func(c *gossie.Context) error {
				name := c.GetArg("name")
				columns := c.Args()[1:] // Skip the first argument (table name)

				if len(columns) == 0 {
					return fmt.Errorf("at least one column definition required")
				}

				fmt.Printf("Creating table '%s'...\n", name)
				if c.HasFlag("if-not-exists") {
					fmt.Printf("CREATE TABLE IF NOT EXISTS %s (\n", name)
				} else {
					fmt.Printf("CREATE TABLE %s (\n", name)
				}

				for i, col := range columns {
					fmt.Printf("  %s", col)
					if i < len(columns)-1 {
						fmt.Println(",")
					} else {
						fmt.Println()
					}
				}
				fmt.Println(");")

				fmt.Printf("Table '%s' created successfully\n", name)
				return nil
			})
		})
	})

	// Backup and restore commands
	app.Command("backup", func(cmd *gossie.Command) {
		cmd.Description("Database backup operations")
		cmd.Arg("destination").Required().Description("Backup destination file")
		cmd.Flag("compress", "Compress the backup").Short('c')
		cmd.Flag("verbose", "Show backup progress").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			destination := c.GetArg("destination")

			fmt.Printf("Starting database backup to %s...\n", destination)

			if c.HasFlag("compress") {
				fmt.Println("Using compression...")
			}

			if c.HasFlag("verbose") {
				fmt.Println("Backing up table: users... 100%")
				fmt.Println("Backing up table: products... 100%")
				fmt.Println("Backing up table: orders... 100%")
				fmt.Println("Backing up table: categories... 100%")
			}

			fmt.Printf("✅ Backup completed successfully: %s\n", destination)
			return nil
		})
	})

	app.Command("restore", func(cmd *gossie.Command) {
		cmd.Description("Restore database from backup")
		cmd.Arg("source").Required().Description("Backup file to restore from")
		cmd.Flag("drop-existing", "Drop existing tables before restore")
		cmd.Flag("verbose", "Show restore progress").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			source := c.GetArg("source")

			if c.HasFlag("drop-existing") {
				fmt.Println("Dropping existing tables...")
			}

			fmt.Printf("Restoring database from %s...\n", source)

			if c.HasFlag("verbose") {
				fmt.Println("Restoring table: users... 100%")
				fmt.Println("Restoring table: products... 100%")
				fmt.Println("Restoring table: orders... 100%")
				fmt.Println("Restoring table: categories... 100%")
			}

			fmt.Printf("✅ Database restored successfully from %s\n", source)
			return nil
		})
	})

	// Migration commands
	app.Command("migrate", func(cmd *gossie.Command) {
		cmd.Description("Database migration operations")

		cmd.Command("status", func(subcmd *gossie.Command) {
			subcmd.Description("Show migration status")
			subcmd.Action(func(c *gossie.Context) error {
				fmt.Println("Migration Status:")
				fmt.Println("=================")
				fmt.Println("Applied migrations:")
				fmt.Println("  001_initial_schema - Applied 2024-01-15 10:30:00")
				fmt.Println("  002_add_user_roles - Applied 2024-01-20 14:45:00")
				fmt.Println("  003_add_indexes - Applied 2024-02-01 09:15:00")
				fmt.Println("")
				fmt.Println("Pending migrations:")
				fmt.Println("  004_add_audit_logs - Pending")
				return nil
			})
		})

		cmd.Command("up", func(subcmd *gossie.Command) {
			subcmd.Description("Apply pending migrations")
			subcmd.Arg("count").Description("Number of migrations to apply")
			subcmd.Flag("dry-run", "Show what would be migrated without applying")
			subcmd.Action(func(c *gossie.Context) error {
				count := c.GetArg("count")

				if c.HasFlag("dry-run") {
					fmt.Println("DRY RUN - Would apply migrations:")
					fmt.Println("  004_add_audit_logs")
					return nil
				}

				if count != "" {
					fmt.Printf("Applying %s migrations...\n", count)
				} else {
					fmt.Println("Applying all pending migrations...")
				}

				fmt.Println("Migration applied: 004_add_audit_logs")
				fmt.Println("✅ All migrations applied successfully")
				return nil
			})
		})

		cmd.Command("down", func(subcmd *gossie.Command) {
			subcmd.Description("Rollback migrations")
			subcmd.Arg("count").Description("Number of migrations to rollback")
			subcmd.Flag("dry-run", "Show what would be rolled back without applying")
			subcmd.Action(func(c *gossie.Context) error {
				count := c.GetArg("count")

				if c.HasFlag("dry-run") {
					fmt.Println("DRY RUN - Would rollback migrations:")
					fmt.Println("  003_add_indexes")
					return nil
				}

				if count != "" {
					fmt.Printf("Rolling back %s migrations...\n", count)
				} else {
					fmt.Println("Rolling back last migration...")
				}

				fmt.Println("Migration rolled back: 003_add_indexes")
				fmt.Println("✅ Rollback completed successfully")
				return nil
			})
		})
	})

	app.Run()
}
