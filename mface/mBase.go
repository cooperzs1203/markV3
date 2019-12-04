package mface

type MBase interface {
	Load() error			// 加载
	Start() error			// 启动
	Stop() error			// 停止
	StartEnding() error		// 开始终止
	OfficialEnding() error	// 正式终止
	Reload() error			// 重新加载
}
