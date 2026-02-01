package config

import "time"

// Config 表示 httpclient 的“行为策略”
// 绝不表示状态
type Config struct {
	Timeout       *time.Duration `yaml:"timeout"`
	DefaultStream *bool          `yaml:"default_stream"`
}

func Default() *Config {
	return &Config{}
}

// Merge 合并配置，永远返回新对象，不修改入参
func Merge(base, override *Config) *Config {
	if base == nil && override == nil {
		return nil
	}

	var c Config
	if base != nil {
		c = *base
	}

	if override != nil {
		if override.Timeout != nil {
			c.Timeout = override.Timeout
		}
		if override.DefaultStream != nil {
			c.DefaultStream = override.DefaultStream
		}
	}

	return &c
}

// EffectiveConfig client / middleware 直接使用
type EffectiveConfig struct {
	Timeout       time.Duration
	DefaultStream bool
}

func toEffective(c *Config) EffectiveConfig {
	var ec EffectiveConfig

	if c == nil {
		return ec
	}

	if c.Timeout != nil {
		ec.Timeout = *c.Timeout
	}

	if c.DefaultStream != nil {
		ec.DefaultStream = *c.DefaultStream
	}

	return ec
}

// RequestContext 是 Resolver 唯一能看到的请求信息
type RequestContext struct {
	URL     string
	Service string
	Headers map[string][]string

	// ⭐ request 级 override（可选）
	Override *Config
}
