package config

import "net/url"

type Resolver interface {
	Resolve(ctx RequestContext) EffectiveConfig
	Version() string
}

type managerResolver struct {
	m *Manager
}

func (r *managerResolver) Resolve(ctx RequestContext) EffectiveConfig {
	snap := r.m.snapshot.Load().(*Snapshot)

	var raw *Config

	svc := ctx.Service
	if svc == "" {
		svc = extractService(ctx.URL)
	}

	if svc != "" {
		raw = snap.ServiceResolved[svc]
	}

	if raw == nil {
		raw = snap.Global
	}

	// ⭐ 只剩 request override
	if ctx.Override != nil {
		raw = Merge(raw, ctx.Override)
	}

	return toEffective(raw)
}

func (r *managerResolver) Version() string {
	snap := r.m.snapshot.Load().(*Snapshot)
	return snap.Version
}

func extractService(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Host
}
