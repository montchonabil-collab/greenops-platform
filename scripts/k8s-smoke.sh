#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${NAMESPACE:-greenops}"

kubectl get pods -n "$NAMESPACE" -o wide
kubectl get svc -n "$NAMESPACE"
kubectl get ingress -n "$NAMESPACE"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=gateway -n "$NAMESPACE" --timeout=120s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=reverse-proxy -n "$NAMESPACE" --timeout=120s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=grafana -n "$NAMESPACE" --timeout=120s

echo "Smoke test Kubernetes OK."
