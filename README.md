# GreenOps Platform

Plateforme SaaS pedagogique pour superviser des metriques energetiques avec une architecture Docker microservices.

## Stack

- Frontend React statique servi par Caddy
- Caddy comme reverse proxy
- API Gateway en Go
- Auth service en Go avec JWT HMAC
- Energy service en Go avec PostgreSQL et Redis
- Alerts service en Go avec PostgreSQL et Redis
- PostgreSQL pour la persistance
- Redis pour le cache
- Prometheus pour les metriques
- Grafana pour les tableaux de bord

## Lancement Docker

Depuis la VM Ubuntu :

```bash
cd /home/tamao/projets/greenops-platform
cp .env.example .env
./scripts/build-linux.sh
docker compose up -d --build
docker compose ps
```

Acces :

- Application : http://192.168.242.132:8080
- API Gateway : http://192.168.242.132:8080/api/health
- Prometheus : http://192.168.242.132:9090
- Grafana : http://192.168.242.132:3001

Compte de demo :

- utilisateur : `admin`
- mot de passe : `greenops`

## Commandes utiles

```bash
docker compose logs -f gateway
docker compose logs -f energy-service
docker compose down
docker compose down -v
```

## Phase Kubernetes

Le dossier `k8s/` contient une base de travail pour migrer l'infrastructure Docker vers Kubernetes.
