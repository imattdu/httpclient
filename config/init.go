package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type initConfig struct {
	filePath        string
	defaultSnapshot *Snapshot
}

type InitOption func(*initConfig)

// WithConfigFile 从配置文件初始化
func WithConfigFile(path string) InitOption {
	return func(c *initConfig) {
		c.filePath = path
	}
}

// WithDefaultSnapshot 设置一个默认 Snapshot（兜底）
func WithDefaultSnapshot(snap *Snapshot) InitOption {
	return func(c *initConfig) {
		c.defaultSnapshot = snap
	}
}

// Init 初始化 httpclient 的进程级配置系统。
// 该方法应在应用启动阶段调用一次。
//
// 典型调用位置：
//   - main()
//   - 服务启动入口
//
// 注意：
//   - Init 是幂等的，多次调用会覆盖旧配置
//   - 不会启动后台 goroutine
//   - 不依赖 client
func Init(opts ...InitOption) error {
	cfg := &initConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// 1️⃣ 初始化默认 snapshot（兜底）
	if cfg.defaultSnapshot != nil {
		Apply(cfg.defaultSnapshot)
	}

	// 2️⃣ 从文件加载（可选）
	if cfg.filePath != "" {
		data, err := os.ReadFile(cfg.filePath)
		if err != nil {
			return err
		}

		var snap Snapshot
		if err := yaml.Unmarshal(data, &snap); err != nil {
			return err
		}

		Apply(&snap)
	}

	return nil
}
