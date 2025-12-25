# gox

[![Go Reference](https://pkg.go.dev/badge/github.com/chinayin/gox.svg)](https://pkg.go.dev/github.com/chinayin/gox)
[![Go Report Card](https://goreportcard.com/badge/github.com/chinayin/gox)](https://goreportcard.com/report/github.com/chinayin/gox)
[![CI](https://github.com/chinayin/gox/workflows/CI/badge.svg)](https://github.com/chinayin/gox/actions)

生产级 Go 工具集：配置管理、数据验证、日志记录等核心工具，助力高效开发。

## 功能特性

### 已发布

- **[cli](cli/)** - 命令行输出格式化与适配器支持
  - 专业启动横幅
  - Cobra 适配器自动提取 CLI 参数
  - 可扩展支持其他 CLI 框架
- **[config](config/)** - 配置管理工具
- **[log](log/)** - 日志工具
- **[idgen](idgen/)** - 统一 ID 生成器
- **[validator](validator/)** - 数据验证工具

### 计划中

- **errors** - 增强的错误处理工具
- **strings** - 字符串操作辅助函数
- **slices** - 切片操作工具
- **maps** - Map 操作工具
- **sync** - 并发工具和线程安全数据结构
- **time** - 时间日期工具
- **crypto** - 加密和哈希工具
- **net** - 网络工具（HTTP 客户端、IP 辅助函数）
- **encoding** - 编码/解码工具（JSON、XML）
- **retry** - 重试机制（支持退避策略）
- **cache** - 内存缓存方案

## 安装

```bash
go get github.com/chinayin/gox
```

## 环境要求

- Go 1.25 或更高版本

## 使用方法

```go
import "github.com/chinayin/gox"
```



## 开发指南

### 代码质量检查

本项目使用 [golangci-lint](https://golangci-lint.run/) 进行代码质量检查。

```bash
# 运行代码检查
make lint

# 自动修复问题
make lint-fix

# 运行测试
make test

# 运行所有检查
make check
```



## 文档

完整文档请访问 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox)。

## 贡献

欢迎贡献代码！请随时提交 Pull Request。

## 更新日志

查看 [CHANGELOG.md](./CHANGELOG.md) 了解版本历史。
