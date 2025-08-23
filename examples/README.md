# Gossie CLI Framework Examples

This directory contains comprehensive examples showcasing the capabilities of the Gossie CLI framework. Each example demonstrates different patterns, features, and use cases.

## 🚀 Quick Start

```bash
# Build and run any example
cd examples/[example-name]
go run main.go --help

# Or run directly from the examples directory
go run ./[example-name]/main.go [command] [args...]
```

## 📁 Available Examples

### 1. **Basic** (`basic/`)

The simplest example showing core framework features:

- Traditional command registration
- Struct-based command registration
- Basic argument and flag handling
- Named argument access

```bash
go run ./basic/main.go --help
go run ./basic/main.go hello
go run ./basic/main.go greet "Alice"
go run ./basic/main.go info --verbose
go run ./basic/main.go echo-command "Hello World" --reverse
```

### 2. **Calculator** (`calculator/`)

Advanced mathematical operations with multiple approaches:

- Complex argument parsing
- Multiple operation modes
- Flag combinations
- Statistical functions
- Trigonometric calculations

```bash
go run ./calculator/main.go --help
go run ./calculator/main.go add 10 20 30
go run ./calculator/main.go divide 100 5 --remainder
go run ./calculator/main.go stats mean 1 2 3 4 5
go run ./calculator/main.go sin 45
```

### 3. **Git-like CLI** (`git-like/`)

Git-inspired version control interface:

- Multiple subcommands
- Complex flag handling
- Repository management
- Branch operations
- Remote management

```bash
go run ./git-like/main.go --help
go run ./git-like/main.go init --bare
go run ./git-like/main.go status --short
go run ./git-like/main.go commit --message "Initial commit" --author "John Doe"
go run ./git-like/main.go branch feature-x
go run ./git-like/main.go remote add origin https://github.com/user/repo.git
```

### 4. **Database CLI** (`database/`)

Database management and operations:

- Connection management
- Query execution
- Schema operations
- Backup and restore
- Migration handling

```bash
go run ./database/main.go --help
go run ./database/main.go connect --host localhost --port 5432 --database mydb
go run ./database/main.go query --query "SELECT * FROM users" --format json
go run ./database/main.go schema show users
go run ./database/main.go backup --destination ./backup.sql --compress
```

### 5. **Configuration Manager** (`config-manager/`)

Advanced configuration management:

- Multi-environment configs
- Type validation
- Configuration diff
- Backup and validation
- Hierarchical settings

```bash
go run ./config-manager/main.go --help
go run ./config-manager/main.go set app.name "MyApp" --type string
go run ./config-manager/main.go get --section database
go run ./config-manager/main.go list --section server
go run ./config-manager/main.go validate --environment production
```

### 6. **File Operations** (`file-operations/`)

Comprehensive file manipulation:

- Copy, move, remove operations
- Directory traversal
- File information
- Batch operations
- Recursive operations

```bash
go run ./file-operations/main.go --help
go run ./file-operations/main.go copy --source ./file.txt --destination ./backup.txt --force
go run ./file-operations/main.go info --path ./document.pdf --detailed --human
go run ./file-operations/main.go list --all --long --human
go run ./file-operations/main.go find --name "*.go" --type f
```

### 7. **HTTP Client** (`http-client/`)

Full-featured HTTP client:

- All HTTP methods
- Header management
- JSON formatting
- Authentication helpers
- API testing utilities

```bash
go run ./http-client/main.go --help
go run ./http-client/main.go get --url https://api.example.com/users --json --verbose
go run ./http-client/main.go post --url https://api.example.com/users --data '{"name":"John"}' --json
go run ./http-client/main.go test latency --url https://api.example.com/health --count 10
go run ./http-client/main.go auth basic --username admin --password secret
```

### 8. **Deployment Tool** (`deployment/`)

DevOps deployment management:

- Multi-environment deployment
- Build management
- Rollback capabilities
- Status monitoring
- Configuration validation

```bash
go run ./deployment/main.go --help
go run ./deployment/main.go deploy --environment production --service api --version v1.2.3
go run ./deployment/main.go build --service web --registry docker.io/myorg --no-cache
go run ./deployment/main.go status --environment staging --services --health
go run ./deployment/main.go rollback --environment production --force
```

## 🎯 Framework Features Demonstrated

### **Core Features**

- ✅ **Command Registration**: Traditional and struct-based approaches
- ✅ **Argument Parsing**: Positional and named arguments
- ✅ **Flag Handling**: Boolean, value, and multiple flags
- ✅ **Subcommands**: Nested command hierarchies
- ✅ **Help System**: Auto-generated help with examples
- ✅ **Type Conversion**: Automatic type parsing and validation

### **Advanced Features**

- ✅ **Named Access**: `c.GetArg()`, `c.GetFlag()`, `c.HasFlag()`
- ✅ **Struct Binding**: Automatic CLI-to-struct mapping
- ✅ **Validation**: Required fields, type checking
- ✅ **Error Handling**: Comprehensive error reporting
- ✅ **Context Management**: Rich context with parsed data
- ✅ **Multiple Patterns**: Various CLI design patterns

### **Real-world Scenarios**

- ✅ **API Clients**: HTTP operations with authentication
- ✅ **DevOps Tools**: Deployment, monitoring, configuration
- ✅ **Database Tools**: Connection, queries, migrations
- ✅ **File Managers**: Copy, move, search operations
- ✅ **Configuration Systems**: Multi-environment, validation
- ✅ **Version Control**: Git-like operations and workflows

## 🏗️ Architecture Patterns

### **1. Traditional Approach**

```go
app.Command("greet", func(cmd *gossie.Command) {
    cmd.Description("Greet a person")
    cmd.Arg("name").Required().Description("Name to greet")
    cmd.Action(func(c *gossie.Context) error {
        name := c.GetArg("name")
        fmt.Printf("Hello, %s!\n", name)
        return nil
    })
})
```

### **2. Struct-based Approach**

```go
type GreetCommand struct {
    Name string `gossie:"name"`
    Formal bool `gossie:"formal"`
}

app.AddCommand(func(cmd *GreetCommand) error {
    if cmd.Formal {
        fmt.Printf("Good day, %s!\n", cmd.Name)
    } else {
        fmt.Printf("Hello, %s!\n", cmd.Name)
    }
    return nil
})
```

### **3. Hybrid Approach**

```go
app.Command("complex", func(cmd *gossie.Command) {
    cmd.Description("Complex operation")

    // Subcommands with different approaches
    cmd.Command("simple", func(subcmd *gossie.Command) {
        // Traditional approach for simple operations
    })

    // Struct-based subcommands for complex operations
    app.AddCommand(func(cmd *ComplexCommand) error {
        // Struct-based for complex validation
    })
})
```

## 🔧 Building and Testing

### **Build All Examples**

```bash
# Build individual examples
go build ./examples/basic
go build ./examples/calculator
go build ./examples/git-like
# ... etc

# Build all examples
for dir in examples/*/; do
    echo "Building $dir"
    go build "$dir"
done
```

### **Test All Examples**

```bash
# Test individual examples
go run ./examples/basic/main.go --help
go run ./examples/calculator/main.go add 1 2 3

# Test framework functionality
go test ./pkg/gossie
```

## 📚 Learning Path

### **Beginner**

1. Start with `basic/` - Learn core concepts
2. Try `calculator/` - Understand argument parsing
3. Explore `file-operations/` - Master flag handling

### **Intermediate**

1. Study `git-like/` - Learn subcommand hierarchies
2. Examine `database/` - Complex real-world scenarios
3. Review `config-manager/` - Advanced validation patterns

### **Advanced**

1. Analyze `http-client/` - Network operations and APIs
2. Study `deployment/` - DevOps workflows and patterns
3. Create your own examples using the patterns shown

## 🎨 Best Practices

### **Command Design**

- Use consistent naming conventions
- Provide clear descriptions
- Group related operations in subcommands
- Use appropriate argument types

### **Error Handling**

- Validate inputs early
- Provide helpful error messages
- Use appropriate exit codes
- Handle edge cases gracefully

### **User Experience**

- Auto-generate comprehensive help
- Support both short and long flags
- Provide examples in help text
- Support both interactive and scriptable modes

### **Code Organization**

- Separate concerns with subcommands
- Use struct-based commands for complex operations
- Keep argument parsing logic separate
- Document complex operations

## 📄 License

These examples are part of the Gossie CLI framework and follow the same license terms.
