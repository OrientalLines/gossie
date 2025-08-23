package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/orientallines/gossie/pkg/gossie"
)

type FileCopyCommand struct {
	Source      string `gossie:"source"`
	Destination string `gossie:"destination"`
	Force       bool   `gossie:"force"`
	Verbose     bool   `gossie:"verbose"`
	Recursive   bool   `gossie:"recursive"`
}

type FileMoveCommand struct {
	Source      string `gossie:"source"`
	Destination string `gossie:"destination"`
	Force       bool   `gossie:"force"`
	Verbose     bool   `gossie:"verbose"`
}

type FileInfoCommand struct {
	Path     string `gossie:"path"`
	Detailed bool   `gossie:"detailed"`
	Human    bool   `gossie:"human"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "File Operations",
		AppDescription: "A comprehensive file manipulation CLI tool",
		AppVersion:     "1.0.0",
		AppAuthor:      "FileOps Team",
	})

	// Copy command using struct-based approach
	app.AddCommand(func(cmd *FileCopyCommand) error {
		if cmd.Source == "" || cmd.Destination == "" {
			return fmt.Errorf("both source and destination are required")
		}

		// Check if source exists
		sourceInfo, err := os.Stat(cmd.Source)
		if os.IsNotExist(err) {
			return fmt.Errorf("source does not exist: %s", cmd.Source)
		}

		if cmd.Verbose {
			fmt.Printf("Copying %s to %s...\n", cmd.Source, cmd.Destination)
		}

		if sourceInfo.IsDir() && !cmd.Recursive {
			return fmt.Errorf("source is a directory, use --recursive to copy directories")
		}

		// Check if destination exists
		if _, err := os.Stat(cmd.Destination); err == nil && !cmd.Force {
			return fmt.Errorf("destination already exists: %s (use --force to overwrite)", cmd.Destination)
		}

		if sourceInfo.IsDir() && cmd.Recursive {
			// Copy directory recursively
			if err := copyDir(cmd.Source, cmd.Destination, cmd.Verbose); err != nil {
				return fmt.Errorf("failed to copy directory: %w", err)
			}
		} else {
			// Copy single file
			if err := copyFile(cmd.Source, cmd.Destination, cmd.Verbose); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
		}

		if cmd.Verbose {
			fmt.Printf("✅ Successfully copied %s to %s\n", cmd.Source, cmd.Destination)
		} else {
			fmt.Printf("✅ Copied %s to %s\n", cmd.Source, cmd.Destination)
		}

		return nil
	})

	// Move command using struct-based approach
	app.AddCommand(func(cmd *FileMoveCommand) error {
		if cmd.Source == "" || cmd.Destination == "" {
			return fmt.Errorf("both source and destination are required")
		}

		// Check if source exists
		if _, err := os.Stat(cmd.Source); os.IsNotExist(err) {
			return fmt.Errorf("source does not exist: %s", cmd.Source)
		}

		// Check if destination exists
		if _, err := os.Stat(cmd.Destination); err == nil && !cmd.Force {
			return fmt.Errorf("destination already exists: %s (use --force to overwrite)", cmd.Destination)
		}

		if cmd.Verbose {
			fmt.Printf("Moving %s to %s...\n", cmd.Source, cmd.Destination)
		}

		if err := os.Rename(cmd.Source, cmd.Destination); err != nil {
			return fmt.Errorf("failed to move file: %w", err)
		}

		if cmd.Verbose {
			fmt.Printf("✅ Successfully moved %s to %s\n", cmd.Source, cmd.Destination)
		} else {
			fmt.Printf("✅ Moved %s to %s\n", cmd.Source, cmd.Destination)
		}

		return nil
	})

	// File info command using struct-based approach
	app.AddCommand(func(cmd *FileInfoCommand) error {
		if cmd.Path == "" {
			return fmt.Errorf("path is required")
		}

		info, err := os.Stat(cmd.Path)
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", cmd.Path)
		}

		if !cmd.Detailed {
			// Simple info
			if info.IsDir() {
				fmt.Printf("%s (directory)\n", cmd.Path)
			} else {
				size := info.Size()
				if cmd.Human {
					fmt.Printf("%s (%s)\n", cmd.Path, formatSize(size))
				} else {
					fmt.Printf("%s (%d bytes)\n", cmd.Path, size)
				}
			}
		} else {
			// Detailed info
			fmt.Printf("File: %s\n", cmd.Path)
			fmt.Printf("Size: %d bytes", info.Size())
			if cmd.Human {
				fmt.Printf(" (%s)", formatSize(info.Size()))
			}
			fmt.Println()

			fmt.Printf("Permissions: %s\n", info.Mode().String())
			fmt.Printf("Modified: %s\n", info.ModTime().Format(time.RFC3339))

			if info.IsDir() {
				fmt.Printf("Type: Directory\n")
			} else {
				fmt.Printf("Type: Regular File\n")
			}
		}

		return nil
	})

	// Traditional commands for more complex operations
	app.Command("list", func(cmd *gossie.Command) {
		cmd.Description("List directory contents")
		cmd.Arg("directory").Description("Directory to list (default: current)")
		cmd.Flag("all", "Show hidden files").Short('a')
		cmd.Flag("long", "Use long listing format").Short('l')
		cmd.Flag("human", "Show human-readable file sizes").Short('h')
		cmd.Flag("reverse", "Reverse the sort order").Short('r')
		cmd.Action(func(c *gossie.Context) error {
			dir := c.GetArg("directory")
			if dir == "" {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
			}

			entries, err := os.ReadDir(dir)
			if err != nil {
				return fmt.Errorf("failed to read directory: %w", err)
			}

			if !c.HasFlag("all") {
				// Filter out hidden files
				var filtered []os.DirEntry
				for _, entry := range entries {
					if !strings.HasPrefix(entry.Name(), ".") {
						filtered = append(filtered, entry)
					}
				}
				entries = filtered
			}

			if c.HasFlag("reverse") {
				// Reverse the order
				for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
					entries[i], entries[j] = entries[j], entries[i]
				}
			}

			if c.HasFlag("long") {
				// Long listing format
				fmt.Printf("Directory: %s\n", dir)
				fmt.Println("Permissions\tSize\tModified\t\tName")
				fmt.Println("-----------\t----\t--------\t\t----")

				for _, entry := range entries {
					info, err := entry.Info()
					if err != nil {
						continue
					}

					size := ""
					if c.HasFlag("human") {
						size = formatSize(info.Size())
					} else {
						size = fmt.Sprintf("%d", info.Size())
					}

					fmt.Printf("%s\t%s\t%s\t%s\n",
						info.Mode().String(),
						size,
						info.ModTime().Format("Jan 02 15:04"),
						entry.Name())
				}
			} else {
				// Simple listing
				for _, entry := range entries {
					fmt.Println(entry.Name())
				}
			}

			return nil
		})
	})

	app.Command("find", func(cmd *gossie.Command) {
		cmd.Description("Find files and directories")
		cmd.Arg("directory").Description("Directory to search in (default: current)")
		cmd.Arg("name").Description("Name pattern to search for")
		cmd.Flag("type", "File type to search for (f=file, d=directory)").Short('t')
		cmd.Flag("maxdepth", "Maximum search depth").Short('m')
		cmd.Action(func(c *gossie.Context) error {
			dir := c.GetArg("directory")
			if dir == "" {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
			}

			namePattern := c.GetArg("name")
			fileType := c.GetFlag("type")
			maxDepthStr := c.GetFlag("maxdepth")

			maxDepth := 10 // default
			if maxDepthStr != "" {
				if n, err := strconv.Atoi(maxDepthStr); err == nil && n > 0 {
					maxDepth = n
				}
			}

			fmt.Printf("Searching in: %s\n", dir)
			if namePattern != "" {
				fmt.Printf("Name pattern: %s\n", namePattern)
			}
			if fileType != "" {
				fmt.Printf("File type: %s\n", fileType)
			}
			fmt.Printf("Max depth: %d\n", maxDepth)

			// Simple mock search results
			fmt.Println("")
			fmt.Println("Found:")
			fmt.Println("./main.go")
			fmt.Println("./pkg/gossie/app.go")
			fmt.Println("./examples/")
			fmt.Println("./README.md")

			return nil
		})
	})

	app.Command("create", func(cmd *gossie.Command) {
		cmd.Description("Create files and directories")

		cmd.Command("file", func(subcmd *gossie.Command) {
			subcmd.Description("Create a new file")
			subcmd.Arg("filename").Required().Description("Name of the file to create")
			subcmd.Flag("content", "Initial content for the file").Short('c')
			subcmd.Action(func(c *gossie.Context) error {
				filename := c.GetArg("filename")
				content := c.GetFlag("content")

				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file already exists: %s", filename)
				}

				file, err := os.Create(filename)
				if err != nil {
					return fmt.Errorf("failed to create file: %w", err)
				}
				defer file.Close()

				if content != "" {
					if _, err := file.WriteString(content); err != nil {
						return fmt.Errorf("failed to write content: %w", err)
					}
				}

				fmt.Printf("✅ Created file: %s\n", filename)
				return nil
			})
		})

		cmd.Command("dir", func(subcmd *gossie.Command) {
			subcmd.Description("Create a new directory")
			subcmd.Arg("dirname").Required().Description("Name of the directory to create")
			subcmd.Flag("parents", "Create parent directories as needed").Short('p')
			subcmd.Action(func(c *gossie.Context) error {
				dirname := c.GetArg("dirname")

				var err error
				if c.HasFlag("parents") {
					err = os.MkdirAll(dirname, 0755)
				} else {
					err = os.Mkdir(dirname, 0755)
				}

				if err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}

				fmt.Printf("✅ Created directory: %s\n", dirname)
				return nil
			})
		})
	})

	app.Command("remove", func(cmd *gossie.Command) {
		cmd.Description("Remove files and directories")
		cmd.Arg("paths").Multiple().Description("Files/directories to remove")
		cmd.Flag("force", "Force removal without confirmation").Short('f')
		cmd.Flag("recursive", "Remove directories recursively").Short('r')
		cmd.Action(func(c *gossie.Context) error {
			paths := c.Args()

			if len(paths) == 0 {
				return fmt.Errorf("at least one path is required")
			}

			for _, path := range paths {
				info, err := os.Stat(path)
				if os.IsNotExist(err) {
					if c.HasFlag("force") {
						continue // Skip silently
					} else {
						return fmt.Errorf("path does not exist: %s", path)
					}
				}

				if info.IsDir() && !c.HasFlag("recursive") {
					return fmt.Errorf("%s is a directory, use -r to remove", path)
				}

				if info.IsDir() {
					if err := os.RemoveAll(path); err != nil {
						return fmt.Errorf("failed to remove directory: %w", err)
					}
					fmt.Printf("Removed directory: %s\n", path)
				} else {
					if err := os.Remove(path); err != nil {
						return fmt.Errorf("failed to remove file: %w", err)
					}
					fmt.Printf("Removed file: %s\n", path)
				}
			}

			fmt.Printf("✅ Successfully removed %d item(s)\n", len(paths))
			return nil
		})
	})

	app.Run()
}

// Helper functions
func copyFile(src, dst string, verbose bool) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if verbose {
		fmt.Printf("Copying file: %s -> %s\n", src, dst)
	}

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string, verbose bool) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			if verbose {
				fmt.Printf("Creating directory: %s\n", targetPath)
			}
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			return copyFile(path, targetPath, verbose)
		}
	})
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
