#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NAMESPACE="${NAMESPACE:-greenops}"

kubectl apply -f "$ROOT_DIR/k8s/namespace.yaml"
kubectl create secret generic greenops-secrets \
  --namespace "$NAMESPACE" \
  --from-literal=JWT_SECRET="${JWT_SECRET:-change-me-in-production}" \
  --from-literal=POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-greenops}" \
  --from-literal=GRAFANA_ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD:-greenops}" \
  --dry-run=client \
  -o yaml | kubectl apply -f -

kubectl kustomize --load-restrictor=LoadRestrictionsNone "$ROOT_DIR/k8s" | kubectl apply -f -

kubectl rollout status deployment/postgres -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/redis -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/auth-service -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/energy-service -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/alerts-service -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/gateway -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/reverse-proxy -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/prometheus -n "$NAMESPACE" --timeout=180s
kubectl rollout status deployment/grafana -n "$NAMESPACE" --timeout=180s

kubectl get deploy,svc,ingress,hpa,pvc -n "$NAMESPACE"
