package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/orientallines/gossie/pkg/gossie"
)

type HTTPGetCommand struct {
	URL       string   `gossie:"url"`
	Headers   []string `gossie:"header"`
	Verbose   bool     `gossie:"verbose"`
	JSON      bool     `gossie:"json"`
	Timeout   int      `gossie:"timeout"`
	UserAgent string   `gossie:"user-agent"`
	Follow    bool     `gossie:"follow"`
}

type HTTPPostCommand struct {
	URL         string   `gossie:"url"`
	Data        string   `gossie:"data"`
	File        string   `gossie:"file"`
	ContentType string   `gossie:"content-type"`
	Headers     []string `gossie:"header"`
	Verbose     bool     `gossie:"verbose"`
	JSON        bool     `gossie:"json"`
}

func main() {
	app := gossie.NewApp(gossie.AppConfig{
		AppName:        "HTTP Client",
		AppDescription: "A powerful HTTP client for API testing and debugging",
		AppVersion:     "1.0.0",
		AppAuthor:      "HTTP Client Team",
	})

	// GET request using struct-based approach
	app.AddCommand(func(cmd *HTTPGetCommand) error {
		if cmd.URL == "" {
			return fmt.Errorf("URL is required")
		}

		if cmd.Timeout <= 0 {
			cmd.Timeout = 30 // default 30 seconds
		}

		if cmd.Verbose {
			fmt.Printf("Making GET request to: %s\n", cmd.URL)
			fmt.Printf("Timeout: %d seconds\n", cmd.Timeout)
		}

		client := &http.Client{
			Timeout: time.Duration(cmd.Timeout) * time.Second,
		}

		if !cmd.Follow {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}

		req, err := http.NewRequest("GET", cmd.URL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set user agent
		userAgent := cmd.UserAgent
		if userAgent == "" {
			userAgent = "HTTP-Client/1.0.0"
		}
		req.Header.Set("User-Agent", userAgent)

		// Set custom headers
		for _, header := range cmd.Headers {
			if strings.Contains(header, ":") {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}

		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		duration := time.Since(start)

		if cmd.Verbose {
			fmt.Printf("Response status: %s\n", resp.Status)
			fmt.Printf("Response time: %v\n", duration)
			fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
			fmt.Printf("Content-Length: %s\n", resp.Header.Get("Content-Length"))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if cmd.JSON {
			// Pretty print JSON
			var jsonData interface{}
			if err := json.Unmarshal(body, &jsonData); err == nil {
				prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(string(body))
			}
		} else {
			fmt.Println(string(body))
		}

		if cmd.Verbose {
			fmt.Printf("\nRequest completed in %v\n", duration)
		}

		return nil
	})

	// POST request using struct-based approach
	app.AddCommand(func(cmd *HTTPPostCommand) error {
		if cmd.URL == "" {
			return fmt.Errorf("URL is required")
		}

		var body io.Reader

		// Handle different data sources
		if cmd.File != "" {
			// Read from file (mock implementation)
			fmt.Printf("Reading data from file: %s\n", cmd.File)
			body = strings.NewReader("file content would be here")
		} else if cmd.Data != "" {
			body = strings.NewReader(cmd.Data)
		} else {
			body = strings.NewReader("")
		}

		if cmd.Verbose {
			fmt.Printf("Making POST request to: %s\n", cmd.URL)
			if cmd.Data != "" {
				fmt.Printf("Data: %s\n", cmd.Data)
			}
		}

		req, err := http.NewRequest("POST", cmd.URL, body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set content type
		contentType := cmd.ContentType
		if contentType == "" {
			if cmd.JSON {
				contentType = "application/json"
			} else {
				contentType = "application/x-www-form-urlencoded"
			}
		}
		req.Header.Set("Content-Type", contentType)

		// Set custom headers
		for _, header := range cmd.Headers {
			if strings.Contains(header, ":") {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		if cmd.Verbose {
			fmt.Printf("Response status: %s\n", resp.Status)
		}

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if cmd.JSON {
			var jsonData interface{}
			if err := json.Unmarshal(responseBody, &jsonData); err == nil {
				prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(string(responseBody))
			}
		} else {
			fmt.Println(string(responseBody))
		}

		return nil
	})

	// Traditional commands for additional HTTP methods
	app.Command("put", func(cmd *gossie.Command) {
		cmd.Description("Send a PUT request")
		cmd.Arg("url").Required().Description("URL to send request to")
		cmd.Arg("data").Description("Data to send")
		cmd.Flag("verbose", "Enable verbose output").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			url := c.GetArg("url")
			data := c.GetArg("data")

			if c.HasFlag("verbose") {
				fmt.Printf("PUT %s\n", url)
				if data != "" {
					fmt.Printf("Data: %s\n", data)
				}
			}

			// Mock PUT request
			fmt.Printf("✅ PUT request sent to %s\n", url)
			fmt.Println("Response: 200 OK")
			fmt.Println("Content: Resource updated successfully")

			return nil
		})
	})

	app.Command("delete", func(cmd *gossie.Command) {
		cmd.Description("Send a DELETE request")
		cmd.Arg("url").Required().Description("URL to send request to")
		cmd.Flag("force", "Force deletion without confirmation").Short('f')
		cmd.Flag("verbose", "Enable verbose output").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			url := c.GetArg("url")

			if !c.HasFlag("force") {
				fmt.Printf("⚠️  This will delete the resource at %s\n", url)
				fmt.Println("Use --force to confirm deletion")
				return nil
			}

			if c.HasFlag("verbose") {
				fmt.Printf("DELETE %s\n", url)
			}

			// Mock DELETE request
			fmt.Printf("✅ DELETE request sent to %s\n", url)
			fmt.Println("Response: 204 No Content")
			fmt.Println("Resource deleted successfully")

			return nil
		})
	})

	app.Command("head", func(cmd *gossie.Command) {
		cmd.Description("Send a HEAD request")
		cmd.Arg("url").Required().Description("URL to send request to")
		cmd.Flag("verbose", "Show all headers").Short('v')
		cmd.Action(func(c *gossie.Context) error {
			url := c.GetArg("url")

			// Mock HEAD request
			fmt.Printf("HEAD %s\n", url)
			fmt.Println("Response: 200 OK")

			if c.HasFlag("verbose") {
				fmt.Println("Headers:")
				fmt.Println("  Content-Type: text/html")
				fmt.Println("  Content-Length: 1234")
				fmt.Println("  Last-Modified: Wed, 21 Oct 2024 07:28:00 GMT")
				fmt.Println("  ETag: \"abc123\"")
			} else {
				fmt.Println("Content-Type: text/html")
				fmt.Println("Content-Length: 1234")
			}

			return nil
		})
	})

	// API testing commands
	app.Command("test", func(cmd *gossie.Command) {
		cmd.Description("API testing utilities")

		cmd.Command("latency", func(subcmd *gossie.Command) {
			subcmd.Description("Test API latency")
			subcmd.Arg("url").Required().Description("API endpoint to test")
			subcmd.Arg("count").Description("Number of requests to make")
			subcmd.Action(func(c *gossie.Context) error {
				url := c.GetArg("url")
				count := c.GetArg("count")

				if count == "" {
					count = "5"
				}

				fmt.Printf("Testing latency for: %s\n", url)
				fmt.Printf("Number of requests: %s\n", count)

				// Mock latency test
				fmt.Println("Results:")
				fmt.Println("  Min: 23ms")
				fmt.Println("  Max: 156ms")
				fmt.Println("  Avg: 67ms")
				fmt.Println("  P95: 123ms")

				return nil
			})
		})

		cmd.Command("status", func(subcmd *gossie.Command) {
			subcmd.Description("Check API endpoint status")
			subcmd.Arg("urls").Multiple().Description("API endpoints to check")
			subcmd.Flag("continuous", "Monitor continuously").Short('c')
			cmd.Flag("interval", "Monitoring interval in seconds")
			subcmd.Action(func(c *gossie.Context) error {
				urls := c.Args()

				if len(urls) == 0 {
					return fmt.Errorf("at least one URL is required")
				}

				if c.HasFlag("continuous") {
					interval := "5"
					if intervalFlag := c.GetFlag("interval"); intervalFlag != "" {
						interval = intervalFlag
					}
					fmt.Printf("Monitoring %d endpoint(s) every %s seconds...\n", len(urls), interval)
					fmt.Println("Press Ctrl+C to stop")
					// In a real implementation, this would run continuously
					return nil
				}

				fmt.Printf("Checking %d endpoint(s)...\n", len(urls))
				for _, url := range urls {
					// Mock status check
					fmt.Printf("✅ %s - 200 OK\n", url)
				}

				return nil
			})
		})
	})

	// Authentication helpers
	app.Command("auth", func(cmd *gossie.Command) {
		cmd.Description("Authentication and authorization helpers")

		cmd.Command("basic", func(subcmd *gossie.Command) {
			subcmd.Description("Generate Basic Auth header")
			subcmd.Arg("username").Required().Description("Username")
			subcmd.Arg("password").Required().Description("Password")
			subcmd.Action(func(c *gossie.Context) error {
				username := c.GetArg("username")
				password := c.GetArg("password")

				// In a real implementation, this would encode the credentials
				fmt.Printf("Basic Auth header: Basic %s:%s\n", username, password)
				fmt.Println("Use with: --header 'Authorization: Basic <encoded-credentials>'")

				return nil
			})
		})

		cmd.Command("bearer", func(subcmd *gossie.Command) {
			subcmd.Description("Generate Bearer token header")
			subcmd.Arg("token").Required().Description("Bearer token")
			subcmd.Action(func(c *gossie.Context) error {
				token := c.GetArg("token")

				fmt.Printf("Bearer Auth header: Bearer %s\n", token)
				fmt.Println("Use with: --header 'Authorization: Bearer <token>'")

				return nil
			})
		})
	})

	app.Run()
}
