package log

// Level 日志级别
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Format 日志格式
const (
	FormatJSON    = "json"
	FormatConsole = "console"
)

// Output 输出目标
const (
	OutputStdout      = "stdout"
	OutputStderr      = "stderr"
	OutputFileDefault = "runtime/log/app.log" // 默认日志文件路径
)

// Options 日志配置选项
type Options struct {
	Level  string // 日志级别: debug, info, warn, error
	Format string // 日志格式: json, console
	Output string // 输出目标: stdout, stderr, /path/to/file
}

// DefaultOptions 返回默认配置
func DefaultOptions() Options {
	return Options{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: OutputStdout,
	}
}
