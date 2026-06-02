package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"greenops-platform/services/internal/common"
)

func main() {
	service := common.Env("SERVICE_NAME", "gateway")
	metrics := common.NewMetrics(service)

	routes := map[string]*url.URL{
		"/api/auth/":   mustURL(common.Env("AUTH_URL", "http://auth-service:8080")),
		"/api/energy/": mustURL(common.Env("ENERGY_URL", "http://energy-service:8080")),
		"/api/alerts":  mustURL(common.Env("ALERTS_URL", "http://alerts-service:8080")),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", common.Health(service))
	mux.HandleFunc("/api/health", common.Health(service))
	mux.HandleFunc("/metrics", metrics.Handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for prefix, target := range routes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				proxy(target, "/api").ServeHTTP(w, r)
				return
			}
		}
		common.Error(w, metrics, http.StatusNotFound, "route not found")
	})

	addr := ":8080"
	log.Printf("%s listening on %s", service, addr)
	log.Fatal(http.ListenAndServe(addr, common.WithCORS(metrics.Middleware(mux))))
}

func mustURL(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid upstream URL %q: %v", raw, err)
	}
	return parsed
}

func proxy(target *url.URL, stripPrefix string) http.Handler {
	p := httputil.NewSingleHostReverseProxy(target)
	original := p.Director
	p.Director = func(r *http.Request) {
		original(r)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, stripPrefix)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.Host = target.Host
	}
	return p
}
