# Support de soutenance complet

Ce document retrace le travail realise sur GreenOps Platform depuis la preparation de la VM jusqu'au deploiement Kubernetes final. Il sert de support pour expliquer le projet clairement pendant la presentation.

## 1. Objectif general du projet

GreenOps Platform est une plateforme SaaS pedagogique de supervision energetique. L'objectif est de montrer une application moderne construite en microservices, conteneurisee avec Docker, orchestree avec Kubernetes, observable avec Prometheus et Grafana, et livree avec une documentation et une CI/CD.

Le besoin initial etait de respecter un cahier des charges Docker et Kubernetes. Le projet devait donc prouver plusieurs competences :

- construire une application en plusieurs services ;
- creer des images Docker ;
- orchestrer les services avec Docker Compose ;
- ajouter une base PostgreSQL et un cache Redis ;
- exposer l'application avec un reverse proxy ;
- superviser les services avec Prometheus et Grafana ;
- migrer l'architecture vers Kubernetes ;
- gerer les configurations, secrets, volumes, probes, autoscaling et network policies ;
- publier le code sur GitHub avec une pipeline CI.

## 2. Environnement de travail

Le travail a ete realise avec une VM Ubuntu accessible depuis la machine Windows.

Informations utiles :

- VM Ubuntu : `192.168.242.132`
- Nom de la VM : `safran`
- Utilisateur Linux : `tamao`
- Dossier projet sur Windows : `C:\docker-projets\greenops-platform`
- Dossier projet sur la VM : `/home/tamao/projets/greenops-platform`
- Depot GitHub : https://github.com/montchonabil-collab/greenops-platform

La connexion SSH a ete preparee avec une cle dediee afin de pouvoir copier les fichiers, lancer les commandes et administrer la VM plus facilement.

## 3. Lecture du cahier des charges

Le cahier des charges demandait une application GreenOps avec une architecture microservices. Les grands blocs attendus etaient :

- frontend web ;
- API Gateway ;
- service d'authentification ;
- services metiers ;
- PostgreSQL ;
- Redis ;
- reverse proxy ;
- Prometheus ;
- Grafana ;
- Docker Compose ;
- Kubernetes ;
- CI/CD ;
- documentation.

La strategie choisie a ete de construire une application complete mais raisonnable, facile a demontrer et a expliquer.

## 4. Choix techniques

Les choix techniques ont ete faits pour rester simples, robustes et pedagogiques.

### Langage backend

Les services backend sont ecrits en Go. Go est adapte aux microservices car il produit des binaires autonomes, rapides et faciles a conteneuriser.

Services Go crees :

- `gateway` : point d'entree API ;
- `auth-service` : authentification et generation de token ;
- `energy-service` : donnees energetiques ;
- `alerts-service` : alertes basees sur les donnees energetiques.

### Frontend

Le frontend est une page React statique servie par Caddy. Elle affiche un tableau de bord avec :

- indicateurs energetiques ;
- consommation moyenne ;
- part renouvelable ;
- emissions CO2 ;
- alertes recentes.

### Donnees et cache

PostgreSQL stocke les donnees energetiques et les alertes. Redis sert de cache court terme pour ameliorer les reponses.

### Reverse proxy

Caddy sert le frontend et redirige les appels `/api/*` vers le gateway.

### Observabilite

Prometheus collecte les metriques exposees par les services. Grafana affiche un dashboard provisionne automatiquement : `GreenOps Overview`.

## 5. Architecture generale

Flux principal :

```text
Utilisateur
  -> Caddy reverse proxy
  -> Frontend React
  -> API Gateway
  -> auth-service
  -> energy-service
  -> alerts-service
  -> PostgreSQL / Redis
```

Flux monitoring :

```text
Prometheus
  -> gateway / auth-service / energy-service / alerts-service
  -> Grafana
```

Le gateway centralise l'acces aux microservices. Les services internes ne sont pas exposes directement a l'utilisateur.

## 6. Travail realise sur la VM

La premiere phase a ete de verifier que la VM Ubuntu etait joignable.

Actions realisees :

- connexion SSH a la VM ;
- verification du nom d'hote ;
- verification des droits utilisateur ;
- creation et installation d'une cle SSH ;
- verification de Docker ;
- installation des outils manquants si necessaire ;
- creation du dossier projet dans `/home/tamao/projets`.

Cette etape etait importante pour pouvoir deployer et tester directement dans un environnement proche d'un serveur.

## 7. Creation du projet applicatif

Le projet a ete cree dans `C:\docker-projets\greenops-platform`, puis synchronise vers la VM.

Structure principale :

```text
greenops-platform/
  compose.yaml
  README.md
  docs/
  docker/
  frontend/
  k8s/
  monitoring/
  scripts/
  services/
```

Dossier `services/` :

```text
services/
  gateway/
  auth/
  energy/
  alerts/
  internal/common/
```

Le dossier `internal/common` contient les fonctions partagees :

- healthcheck ;
- reponses JSON ;
- gestion CORS ;
- metriques Prometheus ;
- fonctions JWT.

## 8. Fonctionnement des services

### API Gateway

Le gateway est le point d'entree des APIs. Il recoit les appels venant du reverse proxy puis les redirige vers les services internes.

Exemples :

- `/api/health`
- `/api/auth/login`
- `/api/energy/metrics`
- `/api/energy/summary`
- `/api/alerts`

### Auth service

Le service d'authentification fournit une connexion simple pour la demo. Il genere un token JWT avec une signature HMAC.

Compte de demo :

- utilisateur : `admin`
- mot de passe : `greenops`

### Energy service

Le service energy initialise et expose des donnees de consommation energetique.

Il utilise :

- PostgreSQL pour lire et stocker les donnees ;
- Redis pour mettre en cache le resume energetique.

### Alerts service

Le service alerts analyse les donnees et expose des alertes.

Types d'alertes :

- consommation elevee ;
- emissions CO2 critiques ;
- part renouvelable trop basse.

## 9. Phase Docker

La phase Docker a permis de conteneuriser toute l'application.

Elements produits :

- Dockerfiles pour les services Go ;
- `compose.yaml` ;
- reseaux Docker ;
- volumes Docker ;
- variables d'environnement ;
- monitoring Prometheus/Grafana ;
- reverse proxy Caddy.

### Images applicatives

Les services Go sont compiles en binaires Linux puis copies dans des images minimales.

Script utilise :

```bash
./scripts/build-linux.sh
```

### Docker Compose

Commande de lancement Docker :

```bash
docker compose up -d --build
docker compose ps
```

Services Compose :

- `reverse-proxy`
- `gateway`
- `auth-service`
- `energy-service`
- `alerts-service`
- `postgres`
- `redis`
- `prometheus`
- `grafana`

### Reseaux Docker

Les reseaux separent les responsabilites :

- `edge` : exposition publique via Caddy ;
- `backend` : communication Caddy, gateway et services ;
- `data` : acces PostgreSQL et Redis ;
- `monitoring` : Prometheus et Grafana.

Cette separation montre une logique de securite et d'organisation.

### Volumes Docker

Volumes persistants :

- `postgres_data`
- `redis_data`
- `prometheus_data`
- `grafana_data`

Ils evitent de perdre les donnees au redemarrage des conteneurs.

## 10. Observabilite avec Prometheus et Grafana

Chaque service expose un endpoint `/metrics`.

Metriques exposees :

- `greenops_requests_total`
- `greenops_errors_total`
- `greenops_uptime_seconds`

Prometheus collecte ces metriques via `monitoring/prometheus/prometheus.yml`.

Grafana est provisionne automatiquement avec :

- datasource Prometheus ;
- dashboard `GreenOps Overview`.

Le dashboard montre :

- disponibilite des services ;
- requetes par seconde ;
- erreurs applicatives ;
- uptime des services ;
- compteurs bruts par service.

Acces Grafana :

```text
http://192.168.242.132:3001
```

Compte Grafana :

- utilisateur : `admin`
- mot de passe : `greenops`

## 11. GitHub et CI/CD

Le projet a ete publie sur GitHub :

```text
https://github.com/montchonabil-collab/greenops-platform
```

Une pipeline GitHub Actions a ete ajoutee dans `.github/workflows/ci.yml`.

La CI realise :

- checkout du code ;
- installation de Go ;
- tests Go ;
- build des binaires Linux ;
- validation Docker Compose ;
- installation de kubectl ;
- rendu des manifests Kubernetes avec Kustomize.

La CI est passee en succes apres le commit :

```text
Add Kubernetes deployment phase
```

Interet de la CI : prouver que le projet peut etre verifie automatiquement a chaque push.

## 12. Phase Kubernetes

La deuxieme grande phase a ete la migration vers Kubernetes.

Le dossier `k8s/` contient les manifests suivants :

- `namespace.yaml`
- `configmap.yaml`
- `secrets.example.yaml`
- `postgres.yaml`
- `redis.yaml`
- `apps.yaml`
- `frontend.yaml`
- `monitoring.yaml`
- `ingress.yaml`
- `hpa.yaml`
- `network-policies.yaml`
- `kustomization.yaml`

### Namespace

Tous les objets sont deployes dans le namespace :

```bash
greenops
```

Commande :

```bash
kubectl get pods -n greenops
```

### Deployments

Chaque composant applicatif a un Deployment :

- `gateway`
- `auth-service`
- `energy-service`
- `alerts-service`
- `reverse-proxy`
- `postgres`
- `redis`
- `prometheus`
- `grafana`

Les Deployments permettent a Kubernetes de maintenir le nombre de pods desire.

### Services

Chaque composant a un Service interne Kubernetes.

Exemples :

- `gateway`
- `auth-service`
- `energy-service`
- `alerts-service`
- `postgres`
- `redis`

Grace a ces Services, les pods communiquent par DNS interne :

```text
http://auth-service:8080
http://energy-service:8080
http://alerts-service:8080
postgres:5432
redis:6379
```

### ConfigMaps

Les ConfigMaps stockent la configuration non sensible :

- URLs internes ;
- utilisateur PostgreSQL ;
- nom de base ;
- adresse Redis ;
- utilisateur Grafana ;
- configuration Caddy ;
- configuration Prometheus ;
- provisioning Grafana.

### Secrets

Les Secrets stockent les informations sensibles :

- `JWT_SECRET`
- `POSTGRES_PASSWORD`
- `GRAFANA_ADMIN_PASSWORD`

Le fichier `secrets.example.yaml` sert d'exemple. En deploiement reel, le secret est cree par le script `k8s-deploy.sh`.

### PVC et stockage persistant

Des PersistentVolumeClaims sont declares pour :

- PostgreSQL ;
- Redis ;
- Prometheus ;
- Grafana.

Cela permet de conserver les donnees meme si un pod est recree.

### Probes

Les pods applicatifs ont des readiness probes et liveness probes.

Role :

- readiness probe : indique si le pod peut recevoir du trafic ;
- liveness probe : indique si Kubernetes doit redemarrer le pod.

Cela augmente la fiabilite du deploiement.

### HPA

Des HorizontalPodAutoscalers sont crees pour :

- `gateway`
- `auth-service`
- `energy-service`
- `alerts-service`

Ils permettent d'augmenter le nombre de pods selon l'utilisation CPU.

Commande de verification :

```bash
kubectl get hpa -n greenops
```

### NetworkPolicies

Les NetworkPolicies limitent les communications entrantes.

Exemples :

- le reverse proxy peut parler au gateway ;
- le gateway peut parler aux services API ;
- les services API peuvent parler a PostgreSQL et Redis ;
- Grafana peut parler a Prometheus.

Cela montre une logique de securite reseau dans Kubernetes.

## 13. Installation de Kubernetes sur la VM

Pour deployer reellement le projet, `k3s` a ete installe sur la VM.

Pourquoi k3s ?

- plus leger qu'un cluster Kubernetes classique ;
- adapte a une VM locale ;
- inclut Traefik par defaut ;
- inclut un stockage local simple ;
- compatible avec `kubectl`.

Verification du cluster :

```bash
kubectl get nodes
kubectl get pods -A
```

Le node `safran` est passe en etat `Ready`.

## 14. Probleme rencontre : espace disque

Pendant le deploiement Kubernetes, la VM a manque d'espace disque. Kubernetes a detecte une pression disque :

```text
node.kubernetes.io/disk-pressure
```

Effet observe :

- pods evicted ;
- pods en Pending ;
- pods en Error ;
- Grafana et d'autres composants ne pouvaient pas rester stables.

Cause :

- la VM avait un disque petit ;
- Docker Compose tournait encore ;
- Docker et Kubernetes stockaient chacun des images ;
- les images Prometheus/Grafana/PostgreSQL prennent de la place.

Solution appliquee :

- arret de Docker Compose ;
- nettoyage Docker ;
- redemarrage de k3s ;
- reimport des images applicatives dans containerd ;
- relance des pods Kubernetes.

Commande utile pour verifier l'espace :

```bash
df -h /
```

Conclusion : pour faire tourner Docker Compose et Kubernetes en meme temps, il faut augmenter le disque de la VM, idealement a 30 Go minimum.

## 15. Images Kubernetes et containerd

k3s utilise `containerd`, pas Docker directement, pour executer les pods.

Cela signifie que les images construites avec Docker doivent etre importees dans le runtime Kubernetes.

Etapes :

```bash
./scripts/k8s-build-images.sh
docker save greenops-gateway:latest greenops-auth-service:latest greenops-energy-service:latest greenops-alerts-service:latest -o /tmp/greenops-images.tar
sudo k3s ctr -n k8s.io images import /tmp/greenops-images.tar
```

Sans cette importation, les pods applicatifs peuvent rester en :

```text
ErrImagePull
ImagePullBackOff
```

## 16. Deploiement Kubernetes final

Commande de deploiement :

```bash
./scripts/k8s-deploy.sh
```

Commande de verification :

```bash
kubectl get pods -n greenops
kubectl get deploy -n greenops
kubectl get svc -n greenops
kubectl get ingress -n greenops
kubectl get hpa -n greenops
kubectl get pvc -n greenops
```

Etat final attendu :

```text
alerts-service   2/2
auth-service     2/2
energy-service   2/2
gateway          2/2
grafana          1/1
postgres         1/1
prometheus       1/1
redis            1/1
reverse-proxy    2/2
```

Le deploiement final sur la VM a ete valide.

## 17. Acces a l'application

La demo active tourne maintenant avec Kubernetes.

URLs :

- Application : http://192.168.242.132:8080
- Grafana : http://192.168.242.132:3001
- Prometheus : http://192.168.242.132:9090

Des services `systemd` ont ete ajoutes pour maintenir les port-forwards :

- `greenops-app-portforward.service`
- `greenops-grafana-portforward.service`
- `greenops-prometheus-portforward.service`

Ils permettent de garder les URLs accessibles apres redemarrage de la VM.

Verification :

```bash
systemctl status greenops-app-portforward.service
systemctl status greenops-grafana-portforward.service
systemctl status greenops-prometheus-portforward.service
```

## 18. Verification de la communication entre conteneurs

Dans Kubernetes, les conteneurs tournent dans des pods. Les pods communiquent via les Services Kubernetes et le DNS interne.

Commandes de preuve :

```bash
kubectl exec -n greenops deploy/gateway -- wget -qO- http://auth-service:8080/health
kubectl exec -n greenops deploy/gateway -- wget -qO- http://energy-service:8080/health
kubectl exec -n greenops deploy/gateway -- wget -qO- http://alerts-service:8080/health
```

Depuis la machine Windows, on peut aussi tester le gateway :

```bash
curl http://192.168.242.132:8080/api/health
curl http://192.168.242.132:8080/api/energy/summary
curl http://192.168.242.132:8080/api/alerts
```

Si ces commandes repondent, cela prouve que :

- le reverse proxy parle au gateway ;
- le gateway parle aux services ;
- les services parlent a PostgreSQL et Redis ;
- les donnees remontent jusqu'au frontend.

## 19. Commandes importantes pour la soutenance

### Voir les pods

```bash
kubectl get pods -n greenops
```

### Voir tous les objets

```bash
kubectl get all -n greenops
```

### Voir les services

```bash
kubectl get svc -n greenops
```

### Voir les ingress

```bash
kubectl get ingress -n greenops
```

### Voir le stockage

```bash
kubectl get pvc -n greenops
```

### Voir l'autoscaling

```bash
kubectl get hpa -n greenops
```

### Voir les logs

```bash
kubectl logs -n greenops deploy/gateway
kubectl logs -n greenops deploy/energy-service
kubectl logs -n greenops deploy/grafana
```

### Tester l'application

```bash
curl http://192.168.242.132:8080/api/health
curl http://192.168.242.132:8080/api/energy/summary
curl http://192.168.242.132:8080/api/alerts
```

## 20. Demo resilience

Pour montrer que Kubernetes recree un pod automatiquement :

```bash
kubectl get pods -n greenops
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
kubectl get pods -n greenops -w
```

Explication a donner :

Le Deployment souhaite garder 2 replicas du gateway. Si un pod est supprime, Kubernetes en recree un automatiquement pour revenir a l'etat desire.

## 21. Demo scaling

Pour montrer le scaling manuel :

```bash
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Retour a 2 replicas :

```bash
kubectl scale deployment/gateway -n greenops --replicas=2
```

Explication a donner :

Kubernetes permet d'augmenter rapidement le nombre de pods pour absorber plus de trafic. Le HPA peut aussi ajuster ce nombre automatiquement selon le CPU.

## 22. Difference entre Docker et Kubernetes dans le projet

Docker a servi a construire et conteneuriser les services.

Docker Compose a servi a lancer facilement toute la stack sur une seule machine.

Kubernetes sert maintenant a orchestrer l'application :

- planification des pods ;
- redemarrage automatique ;
- services internes ;
- stockage persistant ;
- ingress ;
- autoscaling ;
- network policies ;
- probes.

Sur la VM actuelle, Docker Compose est arrete volontairement pour economiser l'espace disque. La version active est Kubernetes.

Commandes a retenir :

```bash
docker compose ps
kubectl get pods -n greenops
```

Si Docker Compose est arrete, `docker ps` peut etre vide. Ce n'est pas une erreur : les conteneurs actifs sont geres par `k3s/containerd`.

## 23. Limites et ameliorations possibles

Le socle principal du cahier des charges est couvert.

Ameliorations possibles :

- augmenter le disque de la VM pour faire tourner Docker Compose et Kubernetes en meme temps ;
- ajouter un registre Docker prive ou GitHub Container Registry ;
- ajouter des tests unitaires plus complets ;
- ajouter des tests de charge ;
- ajouter Loki pour les logs ;
- ajouter RabbitMQ ou Kafka pour de l'asynchrone ;
- ajouter GitOps avec Argo CD ;
- ajouter TLS/HTTPS sur l'Ingress ;
- externaliser davantage les secrets.

Ces elements sont utiles pour aller plus loin, mais le projet actuel couvre deja les exigences principales Docker, Kubernetes, observabilite, CI/CD et documentation.

## 24. Plan conseille pour l'oral

Ordre de presentation recommande :

1. Presenter l'objectif GreenOps Platform.
2. Montrer le schema d'architecture dans `docs/architecture.md`.
3. Expliquer les microservices.
4. Expliquer la phase Docker et Docker Compose.
5. Montrer le `compose.yaml`.
6. Montrer l'application dans le navigateur.
7. Montrer Prometheus et Grafana.
8. Expliquer la CI GitHub Actions.
9. Expliquer la migration Kubernetes.
10. Montrer les manifests dans `k8s/`.
11. Lancer `kubectl get pods -n greenops`.
12. Montrer `kubectl get svc -n greenops`.
13. Montrer `kubectl get hpa -n greenops`.
14. Faire une demo de suppression de pod.
15. Conclure avec les limites et ameliorations possibles.

## 25. Conclusion finale

GreenOps Platform est maintenant un projet complet qui demontre une chaine DevOps coherente :

- developpement d'une application microservices ;
- conteneurisation Docker ;
- orchestration Docker Compose ;
- supervision Prometheus et Grafana ;
- publication GitHub ;
- validation CI/CD ;
- migration Kubernetes ;
- deploiement reel sur une VM Ubuntu avec k3s ;
- documentation de mise en production et de soutenance.

Le projet est donc pret pour une demonstration technique et pour une soutenance structuree autour des notions Docker, Kubernetes, observabilite et automatisation.
