# 🎯 Gossie

> **The elegant CLI framework for Go** - Build beautiful command-line applications with zero boilerplate

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/orientallines/gossie)](https://goreportcard.com/report/github.com/orientallines/gossie)

Gossie is a modern, lightweight CLI framework that brings the power of struct-based command definition to Go applications. Say goodbye to manual argument parsing and hello to elegant, type-safe command-line interfaces.

## ✨ Features

- 🚀 **Zero Boilerplate** - Define commands using simple struct tags
- 🔧 **Type-Safe** - Automatic type conversion and validation
- 🎨 **Beautiful APIs** - Traditional and struct-based command registration
- 📝 **Auto Help** - Generated help system with examples
- ⚡ **Lightning Fast** - Minimal overhead, maximum performance
- 🛠️ **Flexible** - Supports complex argument patterns and flags
- 🎯 **Intuitive** - Named argument access with `c.GetArg()`

## 📦 Installation

```bash
go get github.com/orientallines/gossie
```

## 🚀 Quick Start

Get started in 30 seconds:

```go
package main

import (
    "fmt"
    "github.com/orientallines/gossie/pkg/gossie"
)

func main() {
    app := gossie.NewApp(gossie.AppConfig{
        AppName:        "MyApp",
        AppDescription: "My awesome CLI app",
        AppVersion:     "1.0.0",
    })

    // Simple command
    app.Command("hello", func(cmd *gossie.Command) {
        cmd.Description("Say hello")
        cmd.Action(func(c *gossie.Context) error {
            fmt.Println("Hello, World! 🌟")
            return nil
        })
    })

    // Struct-based command with arguments
    type GreetCommand struct {
        Name   string `gossie:"name"`
        Formal bool   `gossie:"formal"`
    }

    app.AddCommand(func(cmd *GreetCommand) error {
        if cmd.Formal {
            fmt.Printf("Good day, %s!\n", cmd.Name)
        } else {
            fmt.Printf("Hey %s! 👋\n", cmd.Name)
        }
        return nil
    })

    app.Run()
}
```

Run it:

```bash
go run main.go hello
go run main.go greet "Alice" --formal
```

## 🎯 Why Gossie?

| Feature | Gossie | Flag Package | Cobra |
|---------|---------|--------------|-------|
| **Struct-based Commands** | ✅ | ❌ | ❌ |
| **Zero Boilerplate** | ✅ | ❌ | ❌ |
| **Type Safety** | ✅ | ⚠️ | ⚠️ |
| **Auto Help Generation** | ✅ | ✅ | ✅ |
| **Named Arguments** | ✅ | ❌ | ⚠️ |
| **Simple API** | ✅ | ✅ | ❌ |

## 📚 Examples

### Traditional Command Registration

```go
app.Command("deploy", func(cmd *gossie.Command) {
    cmd.Description("Deploy application to environment")

    cmd.Arg("environment").Required().Description("Target environment")
    cmd.Flag("force", "Force deployment").Short('f')
    cmd.Flag("verbose", "Verbose output").Short('v')

    cmd.Action(func(c *gossie.Context) error {
        env := c.GetArg("environment")
        force := c.HasFlag("force")

        if force {
            fmt.Printf("Force deploying to %s...\n", env)
        } else {
            fmt.Printf("Deploying to %s...\n", env)
        }
        return nil
    })
})
```

### Struct-based Command Registration

```go
type DeployCommand struct {
    Environment string `gossie:"environment,required"`
    Force       bool   `gossie:"force,short=f"`
    Verbose     bool   `gossie:"verbose,short=v"`
    Config      string `gossie:"config"`
}

app.AddCommand(func(cmd *DeployCommand) error {
    if cmd.Force {
        fmt.Printf("🚀 Force deploying %s to %s\n", cmd.Config, cmd.Environment)
    } else {
        fmt.Printf("📦 Deploying %s to %s\n", cmd.Config, cmd.Environment)
    }
    return nil
})
```

## 🏗️ Architecture

Gossie supports multiple patterns for different use cases:

### 1. **Simple Commands**

Perfect for basic operations with minimal arguments.

### 2. **Complex Workflows**

Struct-based commands with validation and type conversion.

### 3. **Hybrid Applications**

Mix traditional and struct-based commands in the same app.

### 4. **Subcommand Hierarchies**

Organize related commands with nested structures.

## 🎨 Advanced Features

### Custom Help and Validation

```go
type DatabaseCommand struct {
    Host     string `gossie:"host,required"`
    Port     int    `gossie:"port"`
    Database string `gossie:"database,required"`
    SSL      bool   `gossie:"ssl"`
}

app.AddCommand(func(cmd *DatabaseCommand) error {
    // Automatic validation happens before this function
    // All required fields are guaranteed to be set

    if cmd.Port == 0 {
        cmd.Port = 5432 // Default port
    }

    fmt.Printf("Connecting to %s:%d/%s (SSL: %t)\n",
        cmd.Host, cmd.Port, cmd.Database, cmd.SSL)

    return nil
})
```

### Named Access Methods

```go
app.Command("process", func(cmd *gossie.Command) {
    cmd.Arg("input").Required()
    cmd.Flag("output", "Output file").Short('o')

    cmd.Action(func(c *gossie.Context) error {
        // Named access - no indexing required!
        input := c.GetArg("input")
        output := c.GetFlag("output")

        // Check for flags
        if c.HasFlag("verbose") {
            fmt.Println("🔍 Processing in verbose mode")
        }

        return processFile(input, output)
    })
})
```

## 🗂️ Real-world Examples

Explore the [examples directory](./examples/) for comprehensive use cases:

- **Calculator** - Mathematical operations with complex argument parsing
- **Git-like CLI** - Version control interface with subcommands
- **Database Tool** - Connection management and query execution
- **HTTP Client** - Full-featured API client with authentication
- **Deployment Tool** - DevOps workflows and configuration management
- **File Operations** - Copy, move, and batch file operations

## 📖 Documentation

- [Examples](./examples/) - Comprehensive usage examples
- [API Reference](./pkg/gossie/) - Complete API documentation
- [Best Practices](./examples/README.md) - Design patterns and conventions

## 🎯 Use Cases

Gossie is perfect for:

- 🔧 **DevOps Tools** - Build deployment and infrastructure tools
- 📊 **Data Processing** - Create ETL pipelines and data analysis tools
- 🌐 **API Clients** - Build command-line interfaces for web services
- 🗃️ **Database Tools** - Database management and migration utilities
- 📁 **File Managers** - Advanced file operations and organization
- ⚙️ **Configuration Tools** - Environment and config management

## 🚀 Performance

- **Zero dependencies** - Pure Go standard library
- **Minimal memory footprint** - Optimized for performance
- **Fast startup** - Quick command execution
- **Efficient parsing** - Optimized argument processing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ for the Go community**

[⭐ Star us on GitHub](https://github.com/orientallines/gossie) • [🐛 Report Issues](https://github.com/orientallines/gossie/issues) • [📚 Examples](./examples/)

</div>
