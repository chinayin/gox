package cli

// CommandAdapter 命令行框架适配器接口
//
// 该接口用于从不同的命令行框架（如 cobra、flag 等）中提取信息，
// 以便在启动横幅中自动显示命令参数和配置。
type CommandAdapter interface {
	// GetName 获取应用名称
	GetName() string

	// GetVersion 获取应用版本
	GetVersion() string

	// GetFlags 获取所有命令行参数信息
	// 返回 map[string]FlagInfo，key 为参数名
	GetFlags() map[string]FlagInfo
}

// FlagInfo 命令行参数信息
type FlagInfo struct {
	Name         string // 参数名称
	Value        string // 当前值
	DefaultValue string // 默认值
	Usage        string // 参数说明
	Changed      bool   // 是否被用户修改（非默认值）
	Type         string // 参数类型：string/int/bool 等
}
