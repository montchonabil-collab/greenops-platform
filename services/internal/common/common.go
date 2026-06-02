package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Metrics struct {
	Service  string
	Started  time.Time
	mu       sync.Mutex
	Requests int64
	Errors   int64
}

func NewMetrics(service string) *Metrics {
	return &Metrics{Service: service, Started: time.Now()}
}

func (m *Metrics) IncRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests++
}

func (m *Metrics) IncError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.IncRequest()
		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP greenops_requests_total Total HTTP requests.\n")
	fmt.Fprintf(w, "# TYPE greenops_requests_total counter\n")
	fmt.Fprintf(w, "greenops_requests_total{service=%q} %d\n", m.Service, m.Requests)
	fmt.Fprintf(w, "# HELP greenops_errors_total Total application errors.\n")
	fmt.Fprintf(w, "# TYPE greenops_errors_total counter\n")
	fmt.Fprintf(w, "greenops_errors_total{service=%q} %d\n", m.Service, m.Errors)
	fmt.Fprintf(w, "# HELP greenops_uptime_seconds Service uptime.\n")
	fmt.Fprintf(w, "# TYPE greenops_uptime_seconds gauge\n")
	fmt.Fprintf(w, "greenops_uptime_seconds{service=%q} %.0f\n", m.Service, time.Since(m.Started).Seconds())
}

func Env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, m *Metrics, status int, message string) {
	if m != nil {
		m.IncError()
	}
	JSON(w, status, map[string]string{"error": message})
}

func Health(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		JSON(w, http.StatusOK, map[string]string{
			"service": service,
			"status":  "ok",
		})
	}
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SignToken(secret, subject, role string, ttl time.Duration) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	payload := map[string]any{
		"sub":  subject,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(ttl).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding
	unsigned := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(payloadJSON)
	sig := hmacSHA256(secret, unsigned)
	return unsigned + "." + enc.EncodeToString(sig), nil
}

func VerifyToken(secret, token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	unsigned := parts[0] + "." + parts[1]
	expected := hmacSHA256(secret, unsigned)
	actual, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(actual, expected) {
		return nil, errors.New("invalid signature")
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return nil, err
	}

	exp, ok := payload["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, errors.New("expired token")
	}

	return payload, nil
}

func BearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}

func hmacSHA256(secret, value string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	return mac.Sum(nil)
}
