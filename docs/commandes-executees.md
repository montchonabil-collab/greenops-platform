# Commandes executees et annotees

Ce document regroupe les commandes importantes utilisees pendant la creation, le deploiement et la mise en ligne de GreenOps Platform.

Les commandes sont classees par etape. Certaines commandes ont ete simplifiees pour etre relisibles et reutilisables pendant la soutenance.

Important : les mots de passe et secrets ne sont pas ecrits en clair dans ce document. Quand une commande demande un secret, utiliser une valeur adaptee ou une variable d'environnement.

## 1. Verification des outils sur Windows

### Verifier Git

```powershell
& 'C:\Program Files\Git\cmd\git.exe' --version
```

Annotation : verifie que Git est disponible sur Windows pour initialiser le depot, committer et pousser vers GitHub.

### Verifier GitHub CLI

```powershell
& 'C:\Program Files\GitHub CLI\gh.exe' --version
```

Annotation : verifie que GitHub CLI est disponible pour authentifier le compte GitHub et consulter les workflows.

### Verifier Docker

```powershell
docker --version
docker compose version
```

Annotation : verifie que Docker et Docker Compose sont disponibles. Docker Compose a ete utilise pour valider `compose.yaml` cote Windows et lancer la stack cote VM.

### Verifier kubectl

```powershell
where.exe kubectl
kubectl kustomize --load-restrictor=LoadRestrictionsNone k8s
```

Annotation : verifie que `kubectl` est disponible et que les manifests Kubernetes peuvent etre rendus avec Kustomize.

## 2. Connexion a la VM Ubuntu

### Tester la connexion SSH

```powershell
ssh -i C:\Users\ouatt\.ssh\codex_vm_ed25519 `
  -o BatchMode=yes `
  -o ConnectTimeout=10 `
  -o StrictHostKeyChecking=no `
  -o UserKnownHostsFile=NUL `
  tamao@192.168.242.132 'hostname && whoami'
```

Annotation : verifie que la VM est accessible, que la cle SSH fonctionne et que l'utilisateur connecte est le bon.

### Voir les IP de la VM

```bash
hostname -I
ip -br addr
ip route
```

Annotation : permet d'identifier l'IP de gestion `192.168.242.132` et l'IP bridged `10.9.2.192`.

### Verifier l'espace disque de la VM

```bash
df -h /
du -sh /home/tamao/projets
sudo du -sh /var/lib/rancher /var/lib/docker
```

Annotation : controle l'espace disque, tres important car Docker, Kubernetes, Grafana et Prometheus consomment rapidement du stockage.

## 3. Preparation du projet sur Windows

### Aller dans le dossier projet

```powershell
cd C:\docker-projets\greenops-platform
```

Annotation : le projet principal GreenOps est stocke dans ce dossier sur Windows.

### Voir l'etat Git

```powershell
& 'C:\Program Files\Git\cmd\git.exe' status --short --branch
```

Annotation : verifie si le depot contient des changements avant de modifier ou committer.

### Voir les fichiers du projet

```powershell
Get-ChildItem C:\docker-projets\greenops-platform
Get-ChildItem -Recurse C:\docker-projets\greenops-platform\k8s
```

Annotation : liste les fichiers principaux et les manifests Kubernetes.

## 4. Synchronisation vers la VM

### Copier un fichier vers la VM

```powershell
scp -i C:\Users\ouatt\.ssh\codex_vm_ed25519 `
  C:\docker-projets\greenops-platform\compose.yaml `
  tamao@192.168.242.132:/home/tamao/projets/greenops-platform/compose.yaml
```

Annotation : copie un fichier precis de Windows vers le projet sur la VM.

### Copier un dossier vers la VM

```powershell
scp -r -i C:\Users\ouatt\.ssh\codex_vm_ed25519 `
  C:\docker-projets\greenops-platform\k8s `
  tamao@192.168.242.132:/home/tamao/projets/greenops-platform/
```

Annotation : copie tout le dossier Kubernetes vers la VM.

### Creer les dossiers manquants sur la VM

```bash
mkdir -p /home/tamao/projets/greenops-platform/docs
mkdir -p /home/tamao/projets/greenops-platform/scripts
mkdir -p /home/tamao/projets/greenops-platform/monitoring/grafana/dashboards
mkdir -p /home/tamao/projets/greenops-platform/monitoring/grafana/provisioning/dashboards
```

Annotation : prepare les dossiers qui recoivent les scripts, la documentation et les dashboards Grafana.

## 5. Phase Docker Compose

### Construire les binaires Linux

```bash
cd /home/tamao/projets/greenops-platform
./scripts/build-linux.sh
```

Annotation : compile les services Go en binaires Linux pour les images Docker.

### Valider Docker Compose

```bash
docker compose config --quiet
```

Annotation : verifie la syntaxe et la coherence de `compose.yaml`.

### Lancer la stack Docker

```bash
docker compose up -d --build
```

Annotation : construit et lance tous les services Docker Compose.

### Voir les conteneurs Compose

```bash
docker compose ps
```

Annotation : affiche les conteneurs de la stack Docker Compose.

### Voir les logs d'un service

```bash
docker compose logs -f gateway
docker compose logs -f energy-service
docker compose logs -f grafana
```

Annotation : utile pour diagnostiquer les problemes applicatifs ou Grafana.

### Arreter Docker Compose

```bash
docker compose down
```

Annotation : arrete la stack Docker Compose sans supprimer les volumes.

### Nettoyer Docker

```bash
docker system df
docker system prune -af
```

Annotation : affiche puis libere l'espace disque occupe par Docker. Cette commande a ete utile lorsque la VM manquait d'espace.

## 6. Tests applicatifs Docker ou Kubernetes

### Tester le gateway

```bash
curl http://127.0.0.1:8080/api/health
```

Annotation : verifie que l'API Gateway repond.

### Tester le resume energetique

```bash
curl http://127.0.0.1:8080/api/energy/summary
```

Annotation : verifie que le gateway, le service energy, PostgreSQL et Redis fonctionnent ensemble.

### Tester les alertes

```bash
curl http://127.0.0.1:8080/api/alerts
```

Annotation : verifie que le service alerts communique avec les donnees metier.

### Tester Prometheus

```bash
curl http://127.0.0.1:9090/-/ready
```

Annotation : verifie que Prometheus est pret.

### Tester Grafana

```bash
curl http://127.0.0.1:3001/api/health
```

Annotation : verifie que Grafana est disponible.

## 7. Git et GitHub

### Initialiser ou verifier le depot

```powershell
& 'C:\Program Files\Git\cmd\git.exe' status --short --branch
& 'C:\Program Files\Git\cmd\git.exe' remote -v
```

Annotation : verifie la branche, l'etat local et le depot distant.

### Ajouter les fichiers

```powershell
& 'C:\Program Files\Git\cmd\git.exe' add .
```

Annotation : place les fichiers modifies dans le staging Git.

### Verifier le contenu du commit

```powershell
& 'C:\Program Files\Git\cmd\git.exe' diff --cached --stat
& 'C:\Program Files\Git\cmd\git.exe' diff --check
```

Annotation : controle les fichiers qui seront committes et detecte les erreurs de whitespace.

### Creer un commit

```powershell
& 'C:\Program Files\Git\cmd\git.exe' commit -m "Add Kubernetes deployment phase"
```

Annotation : enregistre une etape importante du projet dans l'historique Git.

### Pousser vers GitHub

```powershell
& 'C:\Program Files\Git\cmd\git.exe' push
```

Annotation : envoie le commit sur GitHub.

### Voir les workflows GitHub Actions

```powershell
& 'C:\Program Files\GitHub CLI\gh.exe' run list --repo montchonabil-collab/greenops-platform --limit 5
```

Annotation : verifie si la CI GitHub Actions est en succes.

## 8. Installation et verification Kubernetes avec k3s

### Installer k3s

```bash
curl -sfL https://get.k3s.io -o /tmp/install-k3s.sh
sudo INSTALL_K3S_EXEC="--write-kubeconfig-mode=644" sh /tmp/install-k3s.sh
```

Annotation : installe Kubernetes leger avec k3s. Le kubeconfig est rendu lisible pour faciliter les commandes `kubectl`.

### Verifier le service k3s

```bash
systemctl status k3s --no-pager -l
systemctl is-active k3s
```

Annotation : controle que Kubernetes tourne correctement.

### Voir les nodes

```bash
kubectl get nodes -o wide
```

Annotation : verifie que le node `safran` est en etat `Ready`.

### Voir les pods systeme

```bash
kubectl get pods -A
```

Annotation : verifie les composants internes de Kubernetes, comme CoreDNS, Traefik et metrics-server.

## 9. Images Kubernetes

### Construire les images locales

```bash
cd /home/tamao/projets/greenops-platform
./scripts/k8s-build-images.sh
```

Annotation : construit les images applicatives `greenops-gateway`, `greenops-auth-service`, `greenops-energy-service` et `greenops-alerts-service`.

### Exporter les images Docker

```bash
docker save \
  greenops-gateway:latest \
  greenops-auth-service:latest \
  greenops-energy-service:latest \
  greenops-alerts-service:latest \
  -o /tmp/greenops-images.tar
```

Annotation : regroupe les images dans une archive.

### Importer les images dans containerd k3s

```bash
sudo k3s ctr -n k8s.io images import /tmp/greenops-images.tar
rm -f /tmp/greenops-images.tar
```

Annotation : k3s utilise `containerd`, pas Docker directement. Cette etape rend les images disponibles pour les pods Kubernetes.

## 10. Deploiement Kubernetes

### Rendre les manifests avec Kustomize

```bash
kubectl kustomize --load-restrictor=LoadRestrictionsNone k8s
```

Annotation : genere le YAML final a partir du dossier `k8s/`.

### Deployer GreenOps

```bash
./scripts/k8s-deploy.sh
```

Annotation : cree le namespace, le Secret de demo, applique les manifests Kubernetes et attend les rollouts.

### Smoke test Kubernetes

```bash
./scripts/k8s-smoke.sh
```

Annotation : verifie que les pods principaux sont prets.

### Voir les pods GreenOps

```bash
kubectl get pods -n greenops -o wide
```

Annotation : remplace `docker ps` pour la version Kubernetes. Les "conteneurs" sont geres par Kubernetes sous forme de pods.

### Voir les deployments

```bash
kubectl get deploy -n greenops
```

Annotation : verifie que les replicas attendus sont disponibles.

### Voir services, ingress, HPA et PVC

```bash
kubectl get svc -n greenops
kubectl get ingress -n greenops
kubectl get hpa -n greenops
kubectl get pvc -n greenops
```

Annotation : controle l'exposition interne, l'Ingress, l'autoscaling et les volumes persistants.

## 11. Diagnostic Kubernetes

### Decrire un pod

```bash
kubectl describe pod -n greenops -l app.kubernetes.io/name=grafana
```

Annotation : utile pour comprendre pourquoi un pod est en Pending, Error ou CrashLoopBackOff.

### Lire les logs d'un deployment

```bash
kubectl logs -n greenops deploy/gateway --tail=120
kubectl logs -n greenops deploy/grafana --tail=120
```

Annotation : permet de diagnostiquer les erreurs applicatives ou de provisioning.

### Voir les evenements Kubernetes

```bash
kubectl get events -n greenops --sort-by=.lastTimestamp
```

Annotation : a permis d'identifier les evictions liees a la pression disque.

### Supprimer les pods echoues

```bash
kubectl delete pod -n greenops --field-selector=status.phase=Failed --ignore-not-found=true
```

Annotation : nettoie les pods evicted ou failed pour permettre aux Deployments de repartir proprement.

### Redemarrer k3s

```bash
sudo systemctl restart k3s
```

Annotation : force Kubernetes a recalculer l'etat du node apres nettoyage disque.

## 12. Port-forward Kubernetes

### Exposer l'application localement sur la VM

```bash
kubectl port-forward -n greenops --address 0.0.0.0 svc/reverse-proxy 8080:80
```

Annotation : rend l'application accessible sur `http://IP_VM:8080`.

### Exposer Grafana localement sur la VM

```bash
kubectl port-forward -n greenops --address 0.0.0.0 svc/grafana 3001:3000
```

Annotation : rend Grafana accessible sur `http://IP_VM:3001`.

### Exposer Prometheus localement sur la VM

```bash
kubectl port-forward -n greenops --address 0.0.0.0 svc/prometheus 9090:9090
```

Annotation : rend Prometheus accessible sur `http://IP_VM:9090`.

### Services systemd de port-forward

```bash
systemctl status greenops-app-portforward.service
systemctl status greenops-grafana-portforward.service
systemctl status greenops-prometheus-portforward.service
```

Annotation : ces services gardent les port-forwards actifs apres redemarrage de la VM.

## 13. Verification reseau de la VM

### Voir les interfaces

```bash
ip -br addr
```

Annotation : a permis de confirmer :

```text
ens33 -> 192.168.242.132
ens34 -> 10.9.2.192
```

### Voir les routes

```bash
ip route
```

Annotation : verifie que la route par defaut passe par le reseau de la salle :

```text
default via 10.9.7.254 dev ens34
```

### Verifier les ports ouverts

```bash
sudo ss -ltnp | grep 8080
sudo ss -ltnp | grep 3001
sudo ss -ltnp | grep 9090
```

Annotation : confirme que les services ecoutent sur `0.0.0.0`, donc sur toutes les interfaces de la VM.

### Verifier le firewall Ubuntu

```bash
sudo ufw status
```

Annotation : verifie que le firewall Ubuntu ne bloque pas les acces.

## 14. Acces depuis les autres machines

### URL locale de salle

```text
http://10.9.2.192:8080
```

Annotation : URL utilisable par les machines capables de joindre directement la VM sur le reseau de la salle.

### Tester depuis Windows

```powershell
Invoke-WebRequest -UseBasicParsing http://10.9.2.192:8080
Invoke-RestMethod -Uri http://10.9.2.192:8080/api/health
```

Annotation : confirme que l'application repond depuis Windows.

### Tester depuis un autre PC Windows

```powershell
Test-NetConnection 10.9.2.192 -Port 8080
```

Annotation : si `TcpTestSucceeded` vaut `False`, le reseau de la salle bloque probablement les connexions directes entre postes.

## 15. Tunnels publics temporaires

### Cloudflare Tunnel teste

```bash
cloudflared tunnel --no-autoupdate --url http://127.0.0.1:8080
```

Annotation : Cloudflare Tunnel a ete teste, mais le reseau bloquait les connexions tunnel vers Cloudflare.

### Variante Cloudflare HTTP/2

```bash
cloudflared tunnel --no-autoupdate --protocol http2 --url http://127.0.0.1:8080
```

Annotation : variante testee pour contourner le blocage UDP/QUIC, mais le reseau bloquait aussi cette methode.

### Tunnel localhost.run pour l'application

```bash
ssh -T \
  -o StrictHostKeyChecking=no \
  -o ServerAliveInterval=30 \
  -o ExitOnForwardFailure=yes \
  -R 80:127.0.0.1:8080 \
  nokey@localhost.run
```

Annotation : cree un tunnel public HTTPS vers l'application GreenOps.

### Tunnel localhost.run pour Grafana

```bash
ssh -T \
  -o StrictHostKeyChecking=no \
  -o ServerAliveInterval=30 \
  -o ExitOnForwardFailure=yes \
  -R 80:127.0.0.1:3001 \
  nokey@localhost.run
```

Annotation : cree un tunnel public HTTPS vers Grafana.

### Tunnel localhost.run pour Prometheus

```bash
ssh -T \
  -o StrictHostKeyChecking=no \
  -o ServerAliveInterval=30 \
  -o ExitOnForwardFailure=yes \
  -R 80:127.0.0.1:9090 \
  nokey@localhost.run
```

Annotation : cree un tunnel public HTTPS vers Prometheus.

### Services systemd des tunnels publics

```bash
systemctl status greenops-localhostrun-tunnel.service
systemctl status greenops-grafana-localhostrun-tunnel.service
systemctl status greenops-prometheus-localhostrun-tunnel.service
```

Annotation : verifie que les tunnels publics sont actifs.

## 16. URLs finales obtenues pendant la mise en ligne

### Application

```text
https://af96cceb82a12a.lhr.life
```

Annotation : tunnel public vers GreenOps.

### Grafana

```text
https://489192f8977e02.lhr.life
```

Annotation : tunnel public vers Grafana.

### Prometheus

```text
https://1b3506ef6ba96a.lhr.life
```

Annotation : tunnel public vers Prometheus.

Attention : ces URLs sont temporaires. Elles peuvent changer si la VM redemarre ou si les services de tunnel se reconnectent.

## 17. Verification des URLs publiques

### Tester l'application publique

```powershell
Invoke-WebRequest -UseBasicParsing https://af96cceb82a12a.lhr.life
Invoke-RestMethod -Uri https://af96cceb82a12a.lhr.life/api/health
```

Annotation : verifie que le frontend et l'API GreenOps sont accessibles publiquement.

### Tester Grafana public

```powershell
Invoke-WebRequest -UseBasicParsing https://489192f8977e02.lhr.life
Invoke-RestMethod -Uri https://489192f8977e02.lhr.life/api/health
```

Annotation : verifie que Grafana est accessible publiquement.

### Tester Prometheus public

```powershell
Invoke-WebRequest -UseBasicParsing https://1b3506ef6ba96a.lhr.life/-/ready
```

Annotation : verifie que Prometheus est accessible publiquement.

## 18. Demonstrations utiles

### Communication entre services

```bash
kubectl exec -n greenops deploy/gateway -- wget -qO- http://auth-service:8080/health
kubectl exec -n greenops deploy/gateway -- wget -qO- http://energy-service:8080/health
kubectl exec -n greenops deploy/gateway -- wget -qO- http://alerts-service:8080/health
```

Annotation : prouve que les pods communiquent via les Services Kubernetes et le DNS interne.

### Resilience Kubernetes

```bash
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
kubectl get pods -n greenops -w
```

Annotation : montre que Kubernetes recree automatiquement les pods manquants.

### Scaling manuel

```bash
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
kubectl scale deployment/gateway -n greenops --replicas=2
```

Annotation : montre que Kubernetes peut augmenter ou reduire le nombre de replicas.

## 19. Commandes pour le mini projet separe

### Creer le dossier du mini projet

```powershell
New-Item -ItemType Directory -Path `
  C:\docker-projets\mini-projet-2,`
  C:\docker-projets\mini-projet-2\incoming,`
  C:\docker-projets\mini-projet-2\docs,`
  C:\docker-projets\mini-projet-2\src `
  -Force
```

Annotation : cree un espace separe pour ne pas impacter GreenOps.

### Initialiser Git dans le mini projet

```powershell
cd C:\docker-projets\mini-projet-2
& 'C:\Program Files\Git\cmd\git.exe' init
& 'C:\Program Files\Git\cmd\git.exe' add .
& 'C:\Program Files\Git\cmd\git.exe' commit -m "Initial mini project scaffold"
```

Annotation : cree un depot Git local separe pour le nouveau mini projet.

## 20. Commandes de controle final

### Controle GreenOps

```bash
kubectl get deploy -n greenops
kubectl get pods -n greenops
systemctl is-active k3s
```

Annotation : verifie que l'application Kubernetes est active.

### Controle Git

```powershell
& 'C:\Program Files\Git\cmd\git.exe' status --short --branch
```

Annotation : verifie que le depot local est propre.

### Controle CI

```powershell
& 'C:\Program Files\GitHub CLI\gh.exe' run list --repo montchonabil-collab/greenops-platform --limit 3
```

Annotation : verifie que les derniers workflows GitHub Actions sont en succes.

## 21. Notes de securite

- Ne pas publier les mots de passe dans GitHub.
- Ne pas laisser Grafana et Prometheus publics plus longtemps que necessaire.
- Les URLs `localhost.run` sont temporaires et peuvent changer.
- Pour une vraie production, utiliser un domaine, HTTPS controle, authentification forte et secrets Kubernetes geres proprement.

## 22. Resume

Les commandes ci-dessus couvrent :

- preparation Windows ;
- connexion SSH a la VM ;
- synchronisation des fichiers ;
- Docker Compose ;
- tests applicatifs ;
- Git et GitHub ;
- installation k3s ;
- deploiement Kubernetes ;
- diagnostic reseau ;
- tunnels publics ;
- creation du mini projet separe.

Ce document peut etre utilise comme annexe technique pendant la soutenance.
