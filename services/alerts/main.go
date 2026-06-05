package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"greenops-platform/services/internal/common"
)

type app struct {
	db      *pgxpool.Pool
	cache   *redis.Client
	metrics *common.Metrics
}

type alert struct {
	ID        int       `json:"id"`
	MetricID  int       `json:"metricId"`
	Site      string    `json:"site"`
	Level     string    `json:"level"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type sourceMetric struct {
	ID               int
	Site             string
	ConsumptionKWh   float64
	RenewablePercent float64
	CO2Kg            float64
}

func main() {
	ctx := context.Background()
	service := common.Env("SERVICE_NAME", "alerts-service")
	metrics := common.NewMetrics(service)

	db := connectDB(ctx, common.Env("DATABASE_URL", "postgres://greenops:GreenOpsPgDemo2026V7nR4qT9sL2@postgres:5432/greenops?sslmode=disable"))
	defer db.Close()

	cache := redis.NewClient(&redis.Options{Addr: common.Env("REDIS_ADDR", "redis:6379")})
	defer cache.Close()

	a := &app{db: db, cache: cache, metrics: metrics}
	if err := a.ensureSchema(ctx); err != nil {
		log.Fatalf("schema init failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", common.Health(service))
	mux.HandleFunc("/metrics", metrics.Handler)
	mux.HandleFunc("/alerts", a.handleAlerts)
	mux.HandleFunc("/alerts/evaluate", a.handleEvaluate)

	addr := ":8080"
	log.Printf("%s listening on %s", service, addr)
	log.Fatal(http.ListenAndServe(addr, common.WithCORS(metrics.Middleware(mux))))
}

func connectDB(ctx context.Context, databaseURL string) *pgxpool.Pool {
	for attempt := 1; attempt <= 30; attempt++ {
		pool, err := pgxpool.New(ctx, databaseURL)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				return pool
			}
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatal("postgres is not reachable")
	return nil
}

func (a *app) ensureSchema(ctx context.Context) error {
	_, err := a.db.Exec(ctx, `
CREATE TABLE IF NOT EXISTS alerts (
  id SERIAL PRIMARY KEY,
  metric_id INT NOT NULL,
  site TEXT NOT NULL,
  level TEXT NOT NULL,
  alert_type TEXT NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(metric_id, alert_type)
);
`)
	return err
}

func (a *app) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		common.Error(w, a.metrics, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	created, err := a.evaluate(r.Context())
	if err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "alert evaluation failed")
		return
	}
	common.JSON(w, http.StatusOK, map[string]any{"created": created})
}

func (a *app) handleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, a.metrics, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if _, err := a.evaluate(r.Context()); err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "alert evaluation failed")
		return
	}

	rows, err := a.db.Query(r.Context(), `
SELECT id, metric_id, site, level, alert_type, message, created_at
FROM alerts
ORDER BY created_at DESC, id DESC
LIMIT 30;
`)
	if err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	items := make([]alert, 0)
	for rows.Next() {
		var item alert
		if err := rows.Scan(&item.ID, &item.MetricID, &item.Site, &item.Level, &item.Type, &item.Message, &item.CreatedAt); err != nil {
			common.Error(w, a.metrics, http.StatusInternalServerError, "row scan failed")
			return
		}
		items = append(items, item)
	}

	common.JSON(w, http.StatusOK, items)
}

func (a *app) evaluate(ctx context.Context) (int, error) {
	rows, err := a.db.Query(ctx, `
SELECT id, site, consumption_kwh, renewable_percent, co2_kg
FROM energy_metrics
ORDER BY recorded_at DESC
LIMIT 24;
`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	created := 0
	for rows.Next() {
		var metric sourceMetric
		if err := rows.Scan(&metric.ID, &metric.Site, &metric.ConsumptionKWh, &metric.RenewablePercent, &metric.CO2Kg); err != nil {
			return 0, err
		}
		candidates := candidatesFor(metric)
		for _, item := range candidates {
			tag, err := a.insertAlert(ctx, metric.ID, metric.Site, item.level, item.kind, item.message)
			if err != nil {
				return 0, err
			}
			created += tag
		}
	}

	_ = a.cache.Set(ctx, "alerts:last_evaluation", time.Now().Format(time.RFC3339), 5*time.Minute).Err()
	return created, nil
}

type candidate struct {
	level   string
	kind    string
	message string
}

func candidatesFor(metric sourceMetric) []candidate {
	items := make([]candidate, 0)
	if metric.ConsumptionKWh >= 900 {
		items = append(items, candidate{"critical", "consumption", "Consommation tres elevee sur " + metric.Site})
	} else if metric.ConsumptionKWh >= 750 {
		items = append(items, candidate{"warning", "consumption", "Consommation elevee sur " + metric.Site})
	}
	if metric.RenewablePercent < 40 {
		items = append(items, candidate{"warning", "renewable", "Part renouvelable sous le seuil sur " + metric.Site})
	}
	if metric.CO2Kg >= 320 {
		items = append(items, candidate{"critical", "co2", "Emissions CO2 critiques sur " + metric.Site})
	}
	return items
}

func (a *app) insertAlert(ctx context.Context, metricID int, site, level, kind, message string) (int, error) {
	tag, err := a.db.Exec(ctx, `
INSERT INTO alerts (metric_id, site, level, alert_type, message)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (metric_id, alert_type) DO NOTHING;
`, metricID, site, level, kind, message)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
