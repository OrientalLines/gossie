package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/orientallines/gossie/pkg/gossie"
)

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "GitLike",
		AppDescription: "A Git-like version control CLI tool",
		AppVersion:     "1.0.0",
		AppAuthor:      "GitLike Team",
	})

	// Repository initialization
	app.Command("init", func(cmd *gossie.Command) {
		cmd.Description("Initialize a new repository")
		cmd.Arg("directory").Description("Directory to initialize (default: current)")
		cmd.Flag("bare", "Create a bare repository")
		cmd.Action(func(c *gossie.Context) error {
			dir := c.GetArg("directory")
			if dir == "" {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
			}

			if c.HasFlag("bare") {
				fmt.Printf("Initializing bare repository in %s\n", dir)
			} else {
				fmt.Printf("Initializing repository in %s\n", dir)
			}

			// Create .gitlike directory structure
			gitDir := filepath.Join(dir, ".gitlike")
			if err := os.MkdirAll(gitDir, 0755); err != nil {
				return fmt.Errorf("failed to create .gitlike directory: %w", err)
			}

			// Create subdirectories
			dirs := []string{"objects", "refs", "refs/heads", "refs/tags"}
			for _, d := range dirs {
				if err := os.MkdirAll(filepath.Join(gitDir, d), 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", d, err)
				}
			}

			// Create HEAD file
			headFile := filepath.Join(gitDir, "HEAD")
			if err := os.WriteFile(headFile, []byte("ref: refs/heads/main\n"), 0644); err != nil {
				return fmt.Errorf("failed to create HEAD file: %w", err)
			}

			fmt.Printf("Repository initialized successfully in %s\n", dir)
			return nil
		})
	})

	// Status command
	app.Command("status", func(cmd *gossie.Command) {
		cmd.Description("Show working tree status")
		cmd.Flag("short", "Give the output in the short-format").Short('s')
		cmd.Flag("porcelain", "Give the output in a stable, easy-to-parse format")
		cmd.Action(func(c *gossie.Context) error {
			if c.HasFlag("porcelain") {
				fmt.Println("# On branch main")
				fmt.Println("# Changes to be committed:")
				fmt.Println("#   (use \"gitlike reset HEAD <file>...\" to unstage)")
				fmt.Println("#")
				fmt.Println("# Changes not staged for commit:")
				fmt.Println("#   (use \"gitlike add <file>...\" to update what will be committed)")
				fmt.Println("#")
				fmt.Println("# Untracked files:")
				fmt.Println("#   (use \"gitlike add <file>...\" to include in what will be committed)")
			} else if c.HasFlag("short") {
				fmt.Println("## main")
				fmt.Println(" M modified-file.txt")
				fmt.Println("?? untracked-file.txt")
			} else {
				fmt.Println("On branch main")
				fmt.Println("")
				fmt.Println("Changes to be committed:")
				fmt.Println("  (use \"gitlike reset HEAD <file>...\" to unstage)")
				fmt.Println("")
				fmt.Println("        modified:   modified-file.txt")
				fmt.Println("")
				fmt.Println("Changes not staged for commit:")
				fmt.Println("  (use \"gitlike add <file>...\" to update what will be committed)")
				fmt.Println("  (use \"gitlike checkout -- <file>...\" to discard changes in working directory)")
				fmt.Println("")
				fmt.Println("        modified:   modified-file.txt")
				fmt.Println("")
				fmt.Println("Untracked files:")
				fmt.Println("  (use \"gitlike add <file>...\" to include in what will be committed)")
				fmt.Println("")
				fmt.Println("        untracked-file.txt")
			}
			return nil
		})
	})

	// Add command
	app.Command("add", func(cmd *gossie.Command) {
		cmd.Description("Add file contents to the index")
		cmd.Arg("files").Multiple().Description("Files to add")
		cmd.Flag("all", "Add all files").Short('A')
		cmd.Flag("verbose", "Be verbose").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			files := c.Args()

			if c.HasFlag("all") {
				if c.HasFlag("verbose") {
					fmt.Println("Adding all files...")
				}
				fmt.Println("Added all files to index")
			} else if len(files) > 0 {
				for _, file := range files {
					if c.HasFlag("verbose") {
						fmt.Printf("Adding %s...\n", file)
					}
				}
				fmt.Printf("Added %d file(s) to index\n", len(files))
			} else {
				return fmt.Errorf("no files specified and --all not used")
			}

			return nil
		})
	})

	// Commit command with struct-based approach
	type CommitCommand struct {
		Message string `gossie:"message"`
		Author  string `gossie:"author"`
		Date    string `gossie:"date"`
		Amend   bool   `gossie:"amend"`
		DryRun  bool   `gossie:"dry-run"`
	}

	app.AddCommand(func(cmd *CommitCommand) error {
		if cmd.Message == "" && !cmd.Amend {
			return fmt.Errorf("commit message is required (use -m flag)")
		}

		if cmd.DryRun {
			fmt.Println("Dry run - would commit with:")
			fmt.Printf("  Message: %s\n", cmd.Message)
			if cmd.Author != "" {
				fmt.Printf("  Author: %s\n", cmd.Author)
			}
			if cmd.Date != "" {
				fmt.Printf("  Date: %s\n", cmd.Date)
			}
			if cmd.Amend {
				fmt.Println("  Amend: true")
			}
			return nil
		}

		// Simulate commit creation
		commitHash := fmt.Sprintf("%x", time.Now().Unix())

		if cmd.Amend {
			fmt.Printf("Amended commit: %s\n", commitHash)
		} else {
			fmt.Printf("Created commit: %s\n", commitHash)
		}

		return nil
	})

	// Log command
	app.Command("log", func(cmd *gossie.Command) {
		cmd.Description("Show commit logs")
		cmd.Flag("oneline", "Show each commit on a single line").Short('o')
		cmd.Flag("graph", "Draw a text-based graphical representation").Short('g')
		cmd.Flag("decorate", "Print ref names").Short('d')
		cmd.Arg("count").Description("Limit the number of commits")
		cmd.Action(func(c *gossie.Context) error {
			count := 10 // default
			if countStr := c.GetArg("count"); countStr != "" {
				if n, err := strconv.Atoi(countStr); err == nil && n > 0 {
					count = n
				}
			}

			if c.HasFlag("oneline") {
				fmt.Println("a1b2c3d Initial commit")
				fmt.Println("e4f5g6h Add README.md")
				fmt.Println("i7j8k9l Update documentation")
			} else if c.HasFlag("graph") {
				fmt.Println("* a1b2c3d - Initial commit (HEAD -> main)")
				fmt.Println("* e4f5g6h - Add README.md")
				fmt.Println("* i7j8k9l - Update documentation")
			} else {
				for i := 0; i < count && i < 3; i++ {
					switch i {
					case 0:
						fmt.Println("commit a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t")
						fmt.Println("Author: John Doe <john@example.com>")
						fmt.Println("Date:   Wed Dec 18 12:34:56 2024 +0000")
						fmt.Println("")
						fmt.Println("    Initial commit")
						fmt.Println("")
					case 1:
						fmt.Println("commit e4f5g6h7i8j9k0l1m2n3o4p5q6r7s8t9u0v1w2")
						fmt.Println("Author: Jane Smith <jane@example.com>")
						fmt.Println("Date:   Wed Dec 18 13:45:00 2024 +0000")
						fmt.Println("")
						fmt.Println("    Add README.md")
						fmt.Println("")
					case 2:
						fmt.Println("commit i7j8k9l0m1n2o3p4q5r6s7t8u9v0w1x2y3z4a5")
						fmt.Println("Author: Bob Wilson <bob@example.com>")
						fmt.Println("Date:   Wed Dec 18 14:56:12 2024 +0000")
						fmt.Println("")
						fmt.Println("    Update documentation")
						fmt.Println("")
					}
				}
			}

			return nil
		})
	})

	// Branch command
	app.Command("branch", func(cmd *gossie.Command) {
		cmd.Description("List, create, or delete branches")
		cmd.Arg("branch").Description("Branch name")
		cmd.Flag("delete", "Delete branch").Short('d')
		cmd.Flag("list", "List branches").Short('l')
		cmd.Flag("all", "List all branches").Short('a')
		cmd.Action(func(c *gossie.Context) error {
			branchName := c.GetArg("branch")

			if c.HasFlag("delete") {
				if branchName == "" {
					return fmt.Errorf("branch name required for delete operation")
				}
				fmt.Printf("Deleted branch %s\n", branchName)
			} else if c.HasFlag("list") || (branchName == "" && !c.HasFlag("delete")) {
				fmt.Println("  develop")
				fmt.Println("* main")
				fmt.Println("  feature/new-ui")
			} else if branchName != "" {
				fmt.Printf("Created branch %s\n", branchName)
			}

			return nil
		})
	})

	// Remote command
	app.Command("remote", func(cmd *gossie.Command) {
		cmd.Description("Manage remote repositories")

		cmd.Command("add", func(subcmd *gossie.Command) {
			subcmd.Description("Add a remote repository")
			subcmd.Arg("name").Required().Description("Remote name")
			subcmd.Arg("url").Required().Description("Remote URL")
			subcmd.Action(func(c *gossie.Context) error {
				name := c.GetArg("name")
				url := c.GetArg("url")
				fmt.Printf("Added remote %s -> %s\n", name, url)
				return nil
			})
		})

		cmd.Command("remove", func(subcmd *gossie.Command) {
			subcmd.Description("Remove a remote repository")
			subcmd.Arg("name").Required().Description("Remote name")
			subcmd.Action(func(c *gossie.Context) error {
				name := c.GetArg("name")
				fmt.Printf("Removed remote %s\n", name)
				return nil
			})
		})

		cmd.Command("list", func(subcmd *gossie.Command) {
			subcmd.Description("List remote repositories")
			subcmd.Flag("verbose", "Show remote URLs").Short('v')
			subcmd.Action(func(c *gossie.Context) error {
				if c.HasFlag("verbose") {
					fmt.Println("origin\thttps://github.com/user/repo.git (fetch)")
					fmt.Println("origin\thttps://github.com/user/repo.git (push)")
					fmt.Println("upstream\thttps://github.com/upstream/repo.git (fetch)")
					fmt.Println("upstream\thttps://github.com/upstream/repo.git (push)")
				} else {
					fmt.Println("origin")
					fmt.Println("upstream")
				}
				return nil
			})
		})

		cmd.Action(func(c *gossie.Context) error {
			fmt.Println("Manage remote repositories:")
			fmt.Println("  remote add <name> <url>    Add a remote")
			fmt.Println("  remote remove <name>       Remove a remote")
			fmt.Println("  remote list                List remotes")
			return nil
		})
	})

	app.Run()
}
