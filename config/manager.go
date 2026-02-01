package config

import (
	"sync/atomic"
)

type Manager struct {
	// snapshot 保存当前 manager「已生效的配置快照」。
	//
	// - 使用 atomic.Value 实现无锁读、原子切换，适用于高并发读场景
	// - snapshot 中存放的是 *Snapshot，不可变对象（immutable）
	// - 每次配置变更都会构造一个新的 Snapshot 并整体替换
	// - Resolve 只做 Load，不会发生数据竞争
	//
	// 设计原则：
	//   - 写少读多
	//   - 读路径无锁
	//   - 配置以“事实快照”的形式存在
	snapshot atomic.Value // *Snapshot
}

func NewManager() *Manager {
	m := &Manager{}
	m.snapshot.Store(&Snapshot{
		Version:         "init",
		Global:          &Config{},
		Services:        map[string]*Config{},
		ServiceResolved: map[string]*Config{},
	})
	return m
}

// ⭐ 进程级唯一 Manager
var defaultManager = NewManager()

// DefaultManager 对外只暴露“获取”
func DefaultManager() *Manager {
	return defaultManager
}

func (m *Manager) Resolver() Resolver {
	return &managerResolver{m: m}
}
