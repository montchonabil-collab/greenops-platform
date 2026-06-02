# Migration Kubernetes

Ce dossier contient la phase Kubernetes de GreenOps Platform. Les manifests sont organises avec Kustomize pour rester lisibles et reproductibles.

## Contenu

- `namespace.yaml` : namespace `greenops`.
- `configmap.yaml` : URLs internes, configuration applicative et compte Grafana.
- `secrets.example.yaml` : exemple de Secret a remplacer en production.
- `postgres.yaml` et `redis.yaml` : services de donnees avec PVC et probes.
- `apps.yaml` : API Gateway et microservices Go.
- `frontend.yaml` : reverse proxy Caddy et frontend statique.
- `monitoring.yaml` : Prometheus, Grafana et volumes persistants.
- `ingress.yaml` : exposition HTTP via Traefik.
- `hpa.yaml` : autoscaling horizontal des services applicatifs.
- `network-policies.yaml` : segmentation reseau entre edge, API, data et monitoring.

## Prerequis

- Un cluster Kubernetes local ou distant avec un Ingress Controller.
- `kubectl` installe et connecte au cluster.
- Des images applicatives disponibles dans le cluster :
  - `greenops-gateway:latest`
  - `greenops-auth-service:latest`
  - `greenops-energy-service:latest`
  - `greenops-alerts-service:latest`

Sur un cluster local qui partage Docker avec l'hote, lance :

```bash
./scripts/k8s-build-images.sh
```

Sur Minikube, il faut construire dans l'environnement Docker de Minikube :

```bash
eval "$(minikube docker-env)"
./scripts/k8s-build-images.sh
```

## Deploiement

```bash
./scripts/k8s-deploy.sh
./scripts/k8s-smoke.sh
```

Le script `k8s-deploy.sh` cree le namespace, applique un Secret de demo, applique Kustomize puis attend les rollouts principaux.

Pour rendre le fichier Secret explicite avant un deploiement manuel :

```bash
cp k8s/secrets.example.yaml /tmp/greenops-secrets.yaml
kubectl apply -f /tmp/greenops-secrets.yaml
kubectl kustomize --load-restrictor=LoadRestrictionsNone k8s | kubectl apply -f -
```

## Acces en demo locale

Sans DNS local :

```bash
kubectl port-forward -n greenops svc/reverse-proxy 8080:80
kubectl port-forward -n greenops svc/prometheus 9090:9090
kubectl port-forward -n greenops svc/grafana 3001:3000
```

Avec Ingress et hosts locaux :

```text
127.0.0.1 greenops.local
127.0.0.1 prometheus.greenops.local
127.0.0.1 grafana.greenops.local
```

URLs :

- Application : http://greenops.local
- Prometheus : http://prometheus.greenops.local
- Grafana : http://grafana.greenops.local

## Verifications utiles

```bash
kubectl get all -n greenops
kubectl get pvc -n greenops
kubectl get hpa -n greenops
kubectl describe ingress greenops -n greenops
kubectl logs -n greenops deploy/gateway
```

## Demo resilience et scaling

Suppression d'un pod applicatif :

```bash
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
kubectl rollout status deployment/gateway -n greenops
```

Scaling manuel :

```bash
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Retour au pilotage par HPA :

```bash
kubectl scale deployment/gateway -n greenops --replicas=2
kubectl get hpa -n greenops
```

## Nettoyage

```bash
kubectl delete namespace greenops
```
