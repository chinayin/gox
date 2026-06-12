package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// sensitiveKeywords 敏感参数名关键字（不区分大小写），命中的参数值在展示时遮蔽
var sensitiveKeywords = []string{
	"token", "password", "passwd", "secret", "key", "auth", "credential", "dsn",
}

// maskedValue 敏感值的展示占位符
const maskedValue = "******"

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
		if isSensitiveName(name) {
			section.Add(name, maskedValue)
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
	// Banner - 显示应用名称和版本
	fmt.Fprintf(s.writer, "\n%s (%s)\n", s.name, s.version)

	// 显示完整命令行
	if fullCmd := getFullCommand(); fullCmd != "" {
		fmt.Fprintf(s.writer, "Command: %s\n", fullCmd)
	}

	fmt.Fprintf(s.writer, "%s\n", strings.Repeat("-", 80))

	// Sections
	for _, section := range s.sections {
		s.printSection(section)
	}

	// Endpoints
	if len(s.endpoints) > 0 {
		fmt.Fprintf(s.writer, "\nServer Endpoints\n")
		for _, ep := range s.endpoints {
			fmt.Fprintf(s.writer, "  %-20s %s\n", ep.Name+":", ep.URL)
		}
	}

	// Footer
	fmt.Fprintf(s.writer, "%s\n", strings.Repeat("-", 80))
	s.printSuccess("Server started successfully")
	fmt.Fprintf(s.writer, "  Press Ctrl+C to shutdown gracefully\n\n")
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

// getFullCommand 获取完整命令行（敏感参数值已遮蔽）
func getFullCommand() string {
	return strings.Join(maskSensitiveArgs(os.Args), " ")
}

// isSensitiveName 判断参数名是否包含敏感关键字（不区分大小写）
func isSensitiveName(name string) bool {
	lower := strings.ToLower(name)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// maskSensitiveArgs 遮蔽命令行参数中的敏感值，
// 支持 --name=value 与 --name value 两种形式，不修改入参切片。
func maskSensitiveArgs(args []string) []string {
	masked := make([]string, len(args))
	copy(masked, args)

	for i := 0; i < len(masked); i++ {
		arg := masked[i]
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		name, _, hasValue := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		if !isSensitiveName(name) {
			continue
		}

		if hasValue {
			prefix, _, _ := strings.Cut(arg, "=")
			masked[i] = prefix + "=" + maskedValue
			continue
		}
		// --name value 形式：遮蔽下一个非 flag 参数
		if i+1 < len(masked) && !strings.HasPrefix(masked[i+1], "-") {
			masked[i+1] = maskedValue
			i++
		}
	}

	return masked
}
