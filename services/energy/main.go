package main

import (
	"context"
	"encoding/json"
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

type energyMetric struct {
	ID               int       `json:"id"`
	Site             string    `json:"site"`
	ConsumptionKWh   float64   `json:"consumptionKwh"`
	RenewablePercent float64   `json:"renewablePercent"`
	CO2Kg            float64   `json:"co2Kg"`
	RecordedAt       time.Time `json:"recordedAt"`
}

type summary struct {
	Points           int     `json:"points"`
	AvgConsumption   float64 `json:"avgConsumptionKwh"`
	AvgRenewable     float64 `json:"avgRenewablePercent"`
	TotalCO2         float64 `json:"totalCo2Kg"`
	HighestSite      string  `json:"highestSite"`
	HighestValueKWh  float64 `json:"highestValueKwh"`
	CacheDescription string  `json:"cacheDescription"`
}

func main() {
	ctx := context.Background()
	service := common.Env("SERVICE_NAME", "energy-service")
	metrics := common.NewMetrics(service)

	db := connectDB(ctx, common.Env("DATABASE_URL", "postgres://greenops:greenops@postgres:5432/greenops?sslmode=disable"))
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
	mux.HandleFunc("/energy/metrics", a.handleMetrics)
	mux.HandleFunc("/energy/summary", a.handleSummary)
	mux.HandleFunc("/energy/simulate", a.handleSimulate)

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
CREATE TABLE IF NOT EXISTS energy_metrics (
  id SERIAL PRIMARY KEY,
  site TEXT NOT NULL,
  consumption_kwh DOUBLE PRECISION NOT NULL,
  renewable_percent DOUBLE PRECISION NOT NULL,
  co2_kg DOUBLE PRECISION NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO energy_metrics (site, consumption_kwh, renewable_percent, co2_kg, recorded_at)
SELECT site, consumption_kwh, renewable_percent, co2_kg, recorded_at
FROM (
  VALUES
    ('Paris HQ', 682.4, 71.2, 184.2, now() - interval '55 minutes'),
    ('Lyon DC', 914.8, 38.5, 332.9, now() - interval '48 minutes'),
    ('Nantes Lab', 421.7, 82.4, 94.6, now() - interval '42 minutes'),
    ('Marseille Hub', 758.2, 44.1, 271.8, now() - interval '34 minutes'),
    ('Paris HQ', 654.1, 73.0, 174.3, now() - interval '25 minutes'),
    ('Lyon DC', 944.2, 35.7, 356.4, now() - interval '18 minutes'),
    ('Nantes Lab', 398.3, 85.2, 82.1, now() - interval '11 minutes'),
    ('Marseille Hub', 733.6, 49.5, 241.7, now() - interval '5 minutes')
) AS seed(site, consumption_kwh, renewable_percent, co2_kg, recorded_at)
WHERE NOT EXISTS (SELECT 1 FROM energy_metrics);
`)
	return err
}

func (a *app) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, a.metrics, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rows, err := a.db.Query(r.Context(), `
SELECT id, site, consumption_kwh, renewable_percent, co2_kg, recorded_at
FROM energy_metrics
ORDER BY recorded_at DESC
LIMIT 24;
`)
	if err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	items := make([]energyMetric, 0)
	for rows.Next() {
		var item energyMetric
		if err := rows.Scan(&item.ID, &item.Site, &item.ConsumptionKWh, &item.RenewablePercent, &item.CO2Kg, &item.RecordedAt); err != nil {
			common.Error(w, a.metrics, http.StatusInternalServerError, "row scan failed")
			return
		}
		items = append(items, item)
	}

	common.JSON(w, http.StatusOK, items)
}

func (a *app) handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		common.Error(w, a.metrics, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	item, err := a.insertDemoMetric(r.Context())
	if err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "metric simulation failed")
		return
	}

	_ = a.cache.Del(r.Context(), "energy:summary").Err()
	common.JSON(w, http.StatusCreated, item)
}

func (a *app) handleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, a.metrics, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context()
	if cached, err := a.cache.Get(ctx, "energy:summary").Result(); err == nil && cached != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "hit")
		_, _ = w.Write([]byte(cached))
		return
	}

	value, err := a.summary(ctx)
	if err != nil {
		common.Error(w, a.metrics, http.StatusInternalServerError, "summary query failed")
		return
	}

	raw, _ := json.Marshal(value)
	_ = a.cache.Set(ctx, "energy:summary", raw, 30*time.Second).Err()
	common.JSON(w, http.StatusOK, value)
}

type demoProfile struct {
	site        string
	consumption float64
	renewable   float64
}

func (a *app) insertDemoMetric(ctx context.Context) (energyMetric, error) {
	profiles := []demoProfile{
		{"Paris HQ", 660, 72},
		{"Lyon DC", 930, 36},
		{"Nantes Lab", 410, 84},
		{"Marseille Hub", 745, 48},
	}

	var count int
	if err := a.db.QueryRow(ctx, `SELECT COUNT(*)::int FROM energy_metrics;`).Scan(&count); err != nil {
		return energyMetric{}, err
	}

	now := time.Now().UTC()
	profile := profiles[count%len(profiles)]
	cycle := count / len(profiles)
	wave := float64((cycle%9)-4) * 18.5
	seconds := float64(now.Second()%12) * 2.4
	consumption := profile.consumption + wave + seconds
	renewable := clamp(profile.renewable+float64((cycle%7)-3)*2.6, 25, 90)
	co2 := consumption * (1 - renewable/100) * 0.58

	var item energyMetric
	err := a.db.QueryRow(ctx, `
INSERT INTO energy_metrics (site, consumption_kwh, renewable_percent, co2_kg, recorded_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, site, consumption_kwh, renewable_percent, co2_kg, recorded_at;
`, profile.site, consumption, renewable, co2, now).Scan(
		&item.ID,
		&item.Site,
		&item.ConsumptionKWh,
		&item.RenewablePercent,
		&item.CO2Kg,
		&item.RecordedAt,
	)
	return item, err
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (a *app) summary(ctx context.Context) (summary, error) {
	var value summary
	err := a.db.QueryRow(ctx, `
WITH latest AS (
  SELECT *
  FROM energy_metrics
  ORDER BY recorded_at DESC
  LIMIT 24
),
highest AS (
  SELECT site, consumption_kwh
  FROM latest
  ORDER BY consumption_kwh DESC
  LIMIT 1
)
SELECT
  COUNT(*)::int,
  COALESCE(AVG(consumption_kwh), 0)::float8,
  COALESCE(AVG(renewable_percent), 0)::float8,
  COALESCE(SUM(co2_kg), 0)::float8,
  COALESCE((SELECT site FROM highest), ''),
  COALESCE((SELECT consumption_kwh FROM highest), 0)::float8
FROM latest;
`).Scan(&value.Points, &value.AvgConsumption, &value.AvgRenewable, &value.TotalCO2, &value.HighestSite, &value.HighestValueKWh)
	value.CacheDescription = "Redis cache TTL 30s"
	return value, err
}
