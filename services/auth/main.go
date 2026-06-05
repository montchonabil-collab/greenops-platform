package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"greenops-platform/services/internal/common"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	service := common.Env("SERVICE_NAME", "auth-service")
	secret := common.Env("JWT_SECRET", "GreenOpsJwtDemo2026u7N9pQ2sV5xL8rT4mK6c")
	metrics := common.NewMetrics(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", common.Health(service))
	mux.HandleFunc("/metrics", metrics.Handler)
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			common.Error(w, metrics, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var input loginRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			common.Error(w, metrics, http.StatusBadRequest, "invalid JSON body")
			return
		}

		role := ""
		switch {
		case input.Username == "admin" && input.Password == "greenops":
			role = "admin"
		case input.Username == "analyst" && input.Password == "greenops":
			role = "analyst"
		default:
			common.Error(w, metrics, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := common.SignToken(secret, input.Username, role, 8*time.Hour)
		if err != nil {
			common.Error(w, metrics, http.StatusInternalServerError, "token generation failed")
			return
		}

		common.JSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  input.Username,
			"role":  role,
		})
	})
	mux.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {
		claims, err := common.VerifyToken(secret, common.BearerToken(r))
		if err != nil {
			common.Error(w, metrics, http.StatusUnauthorized, "invalid token")
			return
		}
		common.JSON(w, http.StatusOK, claims)
	})

	addr := ":8080"
	log.Printf("%s listening on %s", service, addr)
	log.Fatal(http.ListenAndServe(addr, common.WithCORS(metrics.Middleware(mux))))
}
