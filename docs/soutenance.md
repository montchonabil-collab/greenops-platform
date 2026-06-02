# Support de soutenance

## Objectif du projet

GreenOps Platform est une plateforme SaaS pedagogique de supervision energetique. Le projet montre comment passer d'une application microservices conteneurisee avec Docker Compose a une architecture Kubernetes observable, securisee et scalable.

## Points a presenter

1. Architecture generale
   - Frontend React statique servi par Caddy.
   - API Gateway qui centralise les appels `/api`.
   - Services Go independants : auth, energy et alerts.
   - PostgreSQL pour les donnees, Redis pour le cache.
   - Prometheus et Grafana pour l'observabilite.

2. Phase Docker
   - Images applicatives construites a partir de binaires Go.
   - Compose avec reseaux separes : edge, backend, data, monitoring.
   - Volumes persistants pour PostgreSQL, Redis, Prometheus et Grafana.
   - Healthchecks sur PostgreSQL et Redis.

3. Phase Kubernetes
   - Namespace dedie `greenops`.
   - Deployments et Services pour chaque composant.
   - ConfigMaps et Secrets pour separer configuration et donnees sensibles.
   - PVC pour les services persistants.
   - Readiness et liveness probes.
   - HPA pour les services applicatifs.
   - NetworkPolicies pour limiter les flux entrants.

4. CI/CD
   - Tests Go.
   - Build des binaires Linux.
   - Validation Docker Compose.
   - Rendu et validation syntaxique des manifests Kubernetes.

## Scenario de demo

```bash
docker compose up -d --build
docker compose ps
curl http://localhost:8080/api/health
curl http://localhost:8080/api/energy/summary
```

Ouvrir :

- Application : http://localhost:8080
- Prometheus : http://localhost:9090
- Grafana : http://localhost:3001

Puis montrer :

- Les targets Prometheus.
- Le dashboard Grafana `GreenOps Overview`.
- Les logs d'un service : `docker compose logs -f gateway`.

## Scenario Kubernetes

```bash
./scripts/k8s-build-images.sh
./scripts/k8s-deploy.sh
./scripts/k8s-smoke.sh
kubectl get deploy,svc,ingress,hpa,pvc -n greenops
```

Demo resilience :

```bash
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
kubectl rollout status deployment/gateway -n greenops
```

Demo scaling :

```bash
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

## Conclusion

Le projet couvre les attendus Docker et Kubernetes : conteneurisation, orchestration, configuration, persistance, exposition, supervision, securite reseau, automatisation CI et documentation de deploiement.
