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

## Documentation

- `docs/architecture.md` : schema et explication de l'architecture.
- `docs/soutenance.md` : recapitulatif complet du travail realise, du demarrage de la VM au deploiement Kubernetes.
- `docs/commandes-executees.md` : commandes executees et annotees pendant la realisation du projet.
- `docs/installation-machine-physique.md` : guide A a Z pour refaire le projet sur une machine physique Ubuntu.
- `k8s/README.md` : guide de deploiement et de demo Kubernetes.

## Phase Kubernetes

Le dossier `k8s/` contient les manifests Kubernetes : namespace, deployments, services, ingress, secrets, configmaps, PVC, probes, HPA et network policies.

```bash
./scripts/k8s-build-images.sh
./scripts/k8s-deploy.sh
./scripts/k8s-smoke.sh
```

Consulte `k8s/README.md` pour les commandes de demo, de port-forward et de resilience.
