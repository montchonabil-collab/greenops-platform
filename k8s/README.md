# Migration Kubernetes

Base de travail pour la phase 2.

Objets a produire :

- Namespace `greenops`
- Deployments pour gateway, auth, energy, alerts et frontend/reverse proxy
- Services internes pour chaque microservice
- Ingress pour exposer la plateforme
- ConfigMaps pour les URLs et la configuration applicative
- Secrets pour `JWT_SECRET`, `POSTGRES_PASSWORD` et les identifiants sensibles
- PersistentVolumeClaims pour PostgreSQL, Prometheus et Grafana
- Probes `readinessProbe` et `livenessProbe`
- HorizontalPodAutoscaler pour gateway et services metiers

Le dossier sera complete une fois la phase Docker stabilisee.
