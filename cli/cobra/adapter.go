// Package cobra provides a CommandAdapter implementation for spf13/cobra.
package cobra

import (
	"github.com/chinayin/gox/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Adapter Cobra 命令适配器
type Adapter struct {
	cmd *cobra.Command
}

// NewAdapter 创建 Cobra 适配器
func NewAdapter(cmd *cobra.Command) cli.CommandAdapter {
	return &Adapter{cmd: cmd}
}

// GetName 获取应用名称
func (a *Adapter) GetName() string {
	if a.cmd.Short != "" {
		return a.cmd.Short
	}
	return "App"
}

// GetVersion 获取应用版本
func (a *Adapter) GetVersion() string {
	if a.cmd.Version != "" {
		return a.cmd.Version
	}
	return "unknown"
}

// GetFlags 获取所有命令行参数信息
func (a *Adapter) GetFlags() map[string]cli.FlagInfo {
	flags := make(map[string]cli.FlagInfo)

	a.cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = cli.FlagInfo{
			Name:         f.Name,
			Value:        f.Value.String(),
			DefaultValue: f.DefValue,
			Usage:        f.Usage,
			Changed:      f.Changed,
			Type:         f.Value.Type(),
		}
	})

	return flags
}
