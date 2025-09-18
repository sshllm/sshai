package version

import (
	"fmt"
	"runtime"
	"time"
)

// 版本信息变量，在编译时通过 -ldflags 设置
var (
	Version   = "v0.9.18"                           // 版本号
	GitCommit = "unknown"                           // Git提交哈希
	BuildTime = "unknown"                           // 编译时间
	GoVersion = runtime.Version()                   // Go版本
	Platform  = runtime.GOOS + "/" + runtime.GOARCH // 平台信息
)

// BuildInfo 构建信息结构体
type BuildInfo struct {
	Version   string
	GitCommit string
	BuildTime string
	GoVersion string
	Platform  string
}

// GetBuildInfo 获取构建信息
func GetBuildInfo() *BuildInfo {
	return &BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
		Platform:  Platform,
	}
}

// GetVersionString 获取版本字符串
func GetVersionString() string {
	return fmt.Sprintf("SSHAI %s", Version)
}

// GetFullVersionString 获取完整版本信息
func GetFullVersionString() string {
	info := GetBuildInfo()
	return fmt.Sprintf("SSHAI %s (commit: %s, built: %s, go: %s, platform: %s)",
		info.Version, info.GitCommit[:8], info.BuildTime, info.GoVersion, info.Platform)
}

// FormatBuildTime 格式化构建时间
func FormatBuildTime() string {
	if BuildTime == "unknown" {
		return time.Now().Format("2006-01-02 15:04:05 UTC")
	}

	// 尝试解析构建时间
	if t, err := time.Parse("2006-01-02T15:04:05Z", BuildTime); err == nil {
		return t.Format("2006-01-02 15:04:05 UTC")
	}

	return BuildTime
}
