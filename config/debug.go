package config

import "net/http"

func ConfigVersionHandler(r Resolver) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(r.Version()))
	}
}
