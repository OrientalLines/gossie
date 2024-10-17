package gossie

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func printHeader(text string) {
	fmt.Printf("\n%s%s%s\n%s\n", colorBold, colorGreen, text, strings.Repeat("─", len(text)))
}

func printSubHeader(text string) {
	fmt.Printf("\n%s%s%s%s\n", colorBold, colorYellow, text, colorReset)
}

func printAppHelp(app *App) {
	fmt.Printf("%s%s %s%s\n", colorBold, app.name, app.version, colorReset)
	if app.description != "" {
		fmt.Println(app.description)
	}
	fmt.Printf("By %s\n\n", app.author)

	printHeader("USAGE")
	usage := filepath.Base(os.Args[0]) + " [OPTIONS]"
	if len(app.commands) > 0 {
		usage += " <COMMAND> [ARGS...]"
	}
	fmt.Printf("    %s\n", usage)

	if len(app.commands) > 0 {
		printHeader("COMMANDS")
		maxNameLen := 0
		for name := range app.commands {
			if len(name) > maxNameLen {
				maxNameLen = len(name)
			}
		}
		for name, cmd := range app.commands {
			fmt.Printf("    %s%s%s%-*s%s    %s\n", colorBold, colorBlue, name, maxNameLen-len(name), "", colorReset, cmd.description)
		}
	}

	printHeader("OPTIONS")
	fmt.Printf("    %s%s-h, --help%s       Print help information\n", colorBold, colorBlue, colorReset)
	fmt.Printf("    %s%s-V, --version%s    Print version information\n", colorBold, colorBlue, colorReset)

	fmt.Printf("\nUse \"%s <COMMAND> --help\" for more information about a specific command.\n", filepath.Base(os.Args[0]))
}

func printCommandHelp(cmd *Command) {
	printHeader("USAGE")
	usage := fmt.Sprintf("%s %s", cmd.app.name, cmd.name)
	if len(cmd.subcommands) > 0 {
		usage += " <SUBCOMMAND>"
	}
	if len(cmd.flags) > 0 {
		usage += " [OPTIONS]"
	}
	if len(cmd.args) > 0 {
		usage += " [ARGS...]"
	}
	fmt.Printf("%s\n", usage)

	if cmd.description != "" {
		fmt.Printf("\n%s\n", cmd.description)
	}

	if len(cmd.args) > 0 {
		printSubHeader("Arguments:")
		for _, arg := range cmd.args {
			fmt.Printf("  %s%s%-15s%s %s\n", colorBold, colorBlue, arg.name, colorReset, arg.description)
		}
	}

	if len(cmd.flags) > 0 {
		printSubHeader("Flags:")
		for _, flag := range cmd.flags {
			shortFlag := ""
			if flag.shortName != 0 {
				shortFlag = fmt.Sprintf("-%c, ", flag.shortName)
			}
			fmt.Printf("  %s%s%s--%s%s    %s\n", colorBold, colorBlue, shortFlag, flag.name, colorReset, flag.description)
		}
	}

	if len(cmd.subcommands) > 0 {
		printSubHeader("Subcommands:")
		for name, subCmd := range cmd.subcommands {
			fmt.Printf("  %s%s%s%-15s%s %s\n", colorBold, colorBlue, name, "", colorReset, subCmd.description)
		}
		fmt.Printf("\nUse \"%s %s <SUBCOMMAND> --help\" for more information about a specific subcommand.\n", cmd.app.name, cmd.name)
	}

	printSubHeader("Global Flags:")
	fmt.Printf("  %s%s-h, --help%s    Show help for command\n", colorBold, colorBlue, colorReset)
}

func printUsage(name, usage string) {
	if usage == "" {
		usage = "[OPTIONS] [ARGS...]"
	}
	printHeader("USAGE")
	fmt.Printf("%s %s\n", name, usage)
}
