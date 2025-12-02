// Package cli provides professional command-line output formatting.
//
// This package offers standardized startup banners, configuration summaries,
// and status messages for CLI applications, with automatic parameter extraction
// from command-line frameworks.
//
// # Basic Usage
//
// Create a startup banner with configuration sections:
//
//	import "github.com/chinayin/gox/cli"
//
//	startup := cli.NewStartup("MyApp", "v1.0.0").
//		AddSection(
//			cli.NewSection("Configuration").
//				Add("Port", 8080).
//				Add("Workers", 4),
//		).
//		AddEndpoint("Health", "http://localhost:8080/health")
//
//	startup.Print()
//
// Output:
//
//	MyApp v1.0.0
//	--------------------------------------------------------------------------------
//
//	Configuration
//	  Port:                8080
//	  Workers:             4
//
//	Server Endpoints
//	  Health:              http://localhost:8080/health
//
//	--------------------------------------------------------------------------------
//	✓ Server started successfully
//	  Press Ctrl+C to shutdown gracefully
//
// # Cobra Adapter
//
// Automatically extract application info and parameters from Cobra commands:
//
//	import (
//		"github.com/chinayin/gox/cli"
//		clicobra "github.com/chinayin/gox/cli/cobra"
//		"github.com/spf13/cobra"
//	)
//
//	var rootCmd = &cobra.Command{
//		Use:     "myapp",
//		Version: "1.0.0",
//		RunE:    run,
//	}
//
//	func run(cmd *cobra.Command, args []string) error {
//		adapter := clicobra.NewAdapter(cmd)
//
//		// Automatically extracts name, version, and changed flags
//		cli.NewStartupWithAdapter(adapter).
//			AutoAddFlags("help", "version"). // Exclude these flags
//			AddSection(...).
//			Print()
//
//		return nil
//	}
//
// The Parameters section only shows flags that were changed from their default values.
//
// # Custom Adapters
//
// Implement the CommandAdapter interface for other CLI frameworks:
//
//	type CommandAdapter interface {
//		GetName() string
//		GetVersion() string
//		GetFlags() map[string]FlagInfo
//	}
//
// # Color Support
//
// Color output is automatically detected and can be disabled:
//   - Disabled when NO_COLOR environment variable is set
//   - Disabled when TERM=dumb
//   - Disabled when output is not a terminal
//
// # Testing
//
// Redirect output to a buffer for testing:
//
//	var buf bytes.Buffer
//	startup := cli.NewStartup("TestApp", "v1.0.0").
//		WithWriter(&buf)
//	startup.Print()
//	output := buf.String()
//	// Assert on output
//
// # Features
//
//   - Fluent API for easy configuration
//   - Automatic color detection (respects NO_COLOR)
//   - Command adapter support for auto-extracting CLI parameters
//   - Customizable output writer (useful for testing)
//   - Clean, professional formatting
//   - Minimal dependencies (only golang.org/x/term)
package cli
