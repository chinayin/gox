package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Startup provides a fluent interface for printing startup information.
type Startup struct {
	name      string
	version   string
	writer    io.Writer
	useColor  bool
	adapter   CommandAdapter
	sections  []*Section
	endpoints []Endpoint
}

// NewStartup creates a new startup printer.
func NewStartup(name, version string) *Startup {
	return &Startup{
		name:      name,
		version:   version,
		writer:    os.Stdout,
		useColor:  supportsColor(),
		sections:  []*Section{},
		endpoints: []Endpoint{},
	}
}

// NewStartupWithAdapter creates a new startup printer from a command adapter.
// It automatically extracts the application name and version from the adapter.
func NewStartupWithAdapter(adapter CommandAdapter) *Startup {
	return &Startup{
		name:      adapter.GetName(),
		version:   adapter.GetVersion(),
		writer:    os.Stdout,
		useColor:  supportsColor(),
		adapter:   adapter,
		sections:  []*Section{},
		endpoints: []Endpoint{},
	}
}

// WithWriter sets the output writer.
func (s *Startup) WithWriter(w io.Writer) *Startup {
	s.writer = w
	return s
}

// WithAdapter sets the command adapter (optional).
func (s *Startup) WithAdapter(adapter CommandAdapter) *Startup {
	s.adapter = adapter
	return s
}

// AutoAddFlags automatically adds command line flags section (requires adapter).
// excludeNames: flag names to exclude (e.g., "help", "version")
func (s *Startup) AutoAddFlags(excludeNames ...string) *Startup {
	if s.adapter == nil {
		return s
	}

	flags := s.adapter.GetFlags()
	if len(flags) == 0 {
		return s
	}

	section := NewSection("Parameters")
	excludeMap := make(map[string]bool, len(excludeNames))
	for _, name := range excludeNames {
		excludeMap[name] = true
	}

	for name, info := range flags {
		if excludeMap[name] || !info.Changed {
			continue
		}
		section.Add(name, FormatFlagValue(info))
	}

	if len(section.Items) > 0 {
		s.sections = append(s.sections, section)
	}

	return s
}

// AddSection adds a configuration section.
func (s *Startup) AddSection(section *Section) *Startup {
	s.sections = append(s.sections, section)
	return s
}

// AddEndpoint adds a server endpoint.
func (s *Startup) AddEndpoint(name, url string) *Startup {
	s.endpoints = append(s.endpoints, Endpoint{Name: name, URL: url})
	return s
}

// Print outputs all startup information.
func (s *Startup) Print() {
	// Banner
	fmt.Fprintf(s.writer, "\n%s %s\n", s.name, s.version)
	fmt.Fprintln(s.writer, strings.Repeat("-", 80))

	// Sections
	for _, section := range s.sections {
		s.printSection(section)
	}

	// Endpoints
	if len(s.endpoints) > 0 {
		fmt.Fprintln(s.writer, "\nServer Endpoints")
		for _, ep := range s.endpoints {
			fmt.Fprintf(s.writer, "  %-20s %s\n", ep.Name+":", ep.URL)
		}
	}

	// Footer
	fmt.Fprintln(s.writer, strings.Repeat("-", 80))
	s.printSuccess("Server started successfully")
	fmt.Fprint(s.writer, "  Press Ctrl+C to shutdown gracefully\n\n")
}

// printSection outputs a section.
func (s *Startup) printSection(section *Section) {
	fmt.Fprintf(s.writer, "\n%s\n", section.Title)
	for _, item := range section.Items {
		fmt.Fprintf(s.writer, "  %-20s %v\n", item.Key+":", item.Value)
	}
}

// printSuccess prints a success message with optional color.
func (s *Startup) printSuccess(msg string) {
	if s.useColor {
		fmt.Fprintf(s.writer, "\033[32m✓\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(s.writer, "✓ %s\n", msg)
	}
}

// Section represents a configuration section with key-value pairs.
type Section struct {
	Title string
	Items []Item
}

// Item represents a single configuration item.
type Item struct {
	Key   string
	Value any
}

// NewSection creates a new section with the given title.
func NewSection(title string) *Section {
	return &Section{
		Title: title,
		Items: []Item{},
	}
}

// Add adds an item to the section.
func (s *Section) Add(key string, value any) *Section {
	s.Items = append(s.Items, Item{Key: key, Value: value})
	return s
}

// Endpoint represents a server endpoint.
type Endpoint struct {
	Name string
	URL  string
}

// FormatFlagValue formats flag value for display.
func FormatFlagValue(info FlagInfo) string {
	value := info.Value

	// Special handling for bool type
	if info.Type == "bool" {
		if value == "true" {
			return "enabled"
		}
		return "disabled"
	}

	// Mark if value equals default
	if value == info.DefaultValue {
		return value + " (default)"
	}

	return value
}

// supportsColor checks if the terminal supports color output.
func supportsColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}
