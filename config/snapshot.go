package config

// Snapshot 表示一次完整、生效中的配置视图
// 一旦 Apply，视为只读
type Snapshot struct {
	Version  string             `yaml:"version"`
	Global   *Config            `yaml:"global"`
	Services map[string]*Config `yaml:"services"`

	// ⭐ 预计算好的：global + service
	ServiceResolved map[string]*Config `yaml:"-"`
}

func Apply(snap *Snapshot) {
	defaultManager.Apply(snap)
}

// Apply 生效配置
func (m *Manager) Apply(in *Snapshot) {
	if in == nil {
		return
	}

	// 1. 读取当前 snapshot（可能不存在）
	var base *Snapshot
	if cur := m.snapshot.Load(); cur != nil {
		base = cur.(*Snapshot)
	}

	// 2. 构造新的 snapshot（不可变）
	next := &Snapshot{
		Version:  in.Version,
		Global:   nil,
		Services: make(map[string]*Config),
	}

	// 3. 合并 Global
	switch {
	case in.Global != nil && base != nil:
		next.Global = Merge(base.Global, in.Global)
	case in.Global != nil:
		next.Global = in.Global
	case base != nil:
		next.Global = base.Global
	default:
		next.Global = &Config{}
	}

	// 4. 合并 Services（逐个 service merge）
	if base != nil {
		for svc, sc := range base.Services {
			next.Services[svc] = sc
		}
	}

	for svc, sc := range in.Services {
		if baseSc := next.Services[svc]; baseSc != nil {
			next.Services[svc] = Merge(baseSc, sc)
		} else {
			next.Services[svc] = sc
		}
	}

	// 5. 预计算 ServiceResolved
	resolved := make(map[string]*Config, len(next.Services))
	for svc, sc := range next.Services {
		resolved[svc] = Merge(next.Global, sc)
	}
	next.ServiceResolved = resolved

	// 6. 原子切换
	m.snapshot.Store(next)
}
