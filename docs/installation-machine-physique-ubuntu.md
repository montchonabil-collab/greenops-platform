# Installation complete sur une machine physique Ubuntu

Ce document explique comment refaire GreenOps Platform de A a Z sur une machine physique Ubuntu.

Il est pense pour une installation reproductible pendant une soutenance ou un TP : preparation de la machine, installation de Git, Docker, Docker Compose, recuperation du projet, lancement Docker, verification, installation Kubernetes avec k3s, deploiement Kubernetes, exposition reseau, tests et nettoyage.

Les commandes sont donnees pour Ubuntu 22.04 LTS, Ubuntu 24.04 LTS ou une version Ubuntu compatible avec Docker Engine. Adapter seulement le nom d'utilisateur, l'IP de la machine et les secrets.

Important : ne pas mettre de vrais mots de passe dans un document partage. Les valeurs comme `greenops` sont des valeurs de demo locale.

## 1. Objectif final

A la fin, la machine physique doit pouvoir lancer :

- l'application GreenOps sur `http://IP_DE_LA_MACHINE:8080` avec Docker Compose ;
- Grafana sur `http://IP_DE_LA_MACHINE:3001` ;
- Prometheus sur `http://IP_DE_LA_MACHINE:9090` ;
- le meme projet dans Kubernetes avec k3s ;
- des tests de sante pour prouver que les conteneurs et les pods communiquent.

## 2. Hypotheses

Machine cible :

- Ubuntu installe directement sur un PC physique ;
- acces Internet disponible ;
- utilisateur avec droits `sudo` ;
- minimum conseille : 2 CPU, 4 Go RAM, 30 Go disque libre ;
- reseau local fonctionnel pour que les autres machines puissent acceder a l'application.

Dans les commandes, remplacer :

- `UTILISATEUR` par le nom de l'utilisateur Linux ;
- `IP_DE_LA_MACHINE` par l'adresse IP de la machine physique ;
- `EMAIL_GIT` par l'adresse email Git ;
- `NOM_GIT` par le nom Git ;
- les secrets de demo par des valeurs personnelles si le projet n'est pas seulement en demonstration.

## 3. Plan global

1. Mettre Ubuntu a jour.
2. Installer les outils de base.
3. Installer Git.
4. Installer Docker Engine et Docker Compose.
5. Recuperer le code GreenOps.
6. Configurer le fichier `.env`.
7. Lancer le projet avec Docker Compose.
8. Tester l'application et les services.
9. Ouvrir l'acces au reseau local.
10. Installer Kubernetes avec k3s.
11. Construire les images applicatives.
12. Importer les images dans k3s.
13. Deployer les manifests Kubernetes.
14. Tester Kubernetes.
15. Exposer Kubernetes aux autres machines.
16. Nettoyer ou arreter proprement.

## 4. Verification initiale de la machine

Afficher la version Ubuntu :

```bash
cat /etc/os-release
```

Annotation : confirme la distribution et sa version.

Afficher le nom de la machine :

```bash
hostname
hostnamectl
```

Annotation : utile pour identifier la machine sur le reseau et dans Kubernetes.

Afficher l'adresse IP :

```bash
ip -br addr
hostname -I
ip route
```

Annotation : permet de reperer l'IP LAN a donner aux autres utilisateurs.

Verifier les ressources :

```bash
free -h
df -h /
lsblk
nproc
```

Annotation : Docker et Kubernetes demandent de la RAM, du CPU et surtout du disque.

## 5. Mise a jour Ubuntu

Mettre a jour l'index des paquets :

```bash
sudo apt update
```

Installer les mises a jour :

```bash
sudo apt upgrade -y
```

Redemarrer si le systeme le demande :

```bash
sudo reboot
```

Apres redemarrage, revenir sur la machine et verifier :

```bash
uptime
whoami
```

Annotation : demarre sur une base propre avant d'installer Docker et Kubernetes.

## 6. Installation des outils de base

Installer les utilitaires utiles :

```bash
sudo apt install -y \
  ca-certificates \
  curl \
  gnupg \
  lsb-release \
  git \
  openssh-client \
  openssh-server \
  ufw \
  jq \
  nano \
  unzip
```

Annotation : ces outils servent a recuperer le projet, installer Docker, faire des appels HTTP, lire du JSON et administrer la machine.

Activer SSH si la machine doit etre administree a distance :

```bash
sudo systemctl enable --now ssh
sudo systemctl status ssh --no-pager
```

Annotation : SSH permet de se connecter depuis un autre PC avec `ssh UTILISATEUR@IP_DE_LA_MACHINE`.

## 7. Configuration Git

Verifier Git :

```bash
git --version
```

Configurer l'identite Git :

```bash
git config --global user.name "NOM_GIT"
git config --global user.email "EMAIL_GIT"
git config --global init.defaultBranch main
```

Verifier la configuration :

```bash
git config --global --list
```

Annotation : Git sert a recuperer le projet et a garder une trace des modifications.

## 8. Installation Docker Engine

Les commandes ci-dessous suivent la methode officielle Docker avec depot `apt`.

Supprimer les anciens paquets qui peuvent entrer en conflit :

```bash
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
  sudo apt remove -y "$pkg"
done
```

Annotation : evite les conflits entre les paquets Ubuntu anciens et Docker Engine officiel.

Preparer le depot Docker :

```bash
sudo apt update
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
```

Ajouter la source `apt` Docker :

```bash
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF
```

Mettre a jour l'index :

```bash
sudo apt update
```

Installer Docker :

```bash
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Activer Docker au demarrage :

```bash
sudo systemctl enable --now docker
sudo systemctl status docker --no-pager
```

Tester Docker avec `sudo` :

```bash
sudo docker run hello-world
```

Autoriser l'utilisateur courant a utiliser Docker sans `sudo` :

```bash
sudo usermod -aG docker "$USER"
newgrp docker
```

Verifier sans `sudo` :

```bash
docker version
docker compose version
docker run --rm hello-world
```

Annotation : si la commande Docker sans `sudo` echoue, se deconnecter puis se reconnecter.

## 9. Pare-feu local

Autoriser SSH :

```bash
sudo ufw allow OpenSSH
```

Autoriser les ports utiles Docker Compose :

```bash
sudo ufw allow 8080/tcp
sudo ufw allow 3001/tcp
sudo ufw allow 9090/tcp
```

Autoriser HTTP pour Kubernetes Ingress :

```bash
sudo ufw allow 80/tcp
```

Activer le pare-feu :

```bash
sudo ufw enable
sudo ufw status verbose
```

Annotation : Docker publie aussi ses propres regles reseau. Le pare-feu aide pour SSH et l'acces general, mais il faut toujours verifier les ports reellement ouverts.

## 10. Recuperation du projet depuis GitHub

Creer un dossier de projets :

```bash
mkdir -p "$HOME/projets"
cd "$HOME/projets"
```

Cloner le depot :

```bash
git clone https://github.com/montchonabil-collab/greenops-platform.git
cd greenops-platform
```

Verifier l'etat :

```bash
git status
git log --oneline --decorate -5
```

Annotation : le projet est maintenant disponible dans `~/projets/greenops-platform`.

Alternative si la machine physique n'a pas acces a GitHub mais qu'un autre PC possede le dossier :

```bash
scp -r greenops-platform UTILISATEUR@IP_DE_LA_MACHINE:/home/UTILISATEUR/projets/
```

Puis, sur la machine physique :

```bash
cd "$HOME/projets/greenops-platform"
git status
```

## 11. Verification de la structure du projet

Afficher les fichiers principaux :

```bash
ls -la
find . -maxdepth 2 -type f | sort
```

Verifier les scripts :

```bash
ls -la scripts
```

Rendre les scripts executables si besoin :

```bash
chmod +x scripts/*.sh
```

Verifier le fichier Docker Compose :

```bash
docker compose config
```

Annotation : `docker compose config` detecte les erreurs YAML et affiche la configuration finale.

## 12. Configuration des variables d'environnement

Copier l'exemple :

```bash
cp .env.example .env
```

Ouvrir le fichier :

```bash
nano .env
```

Exemple de contenu pour une demo locale :

```text
PROJECT_NAME=greenops
JWT_SECRET=change-me-in-production
POSTGRES_USER=greenops
POSTGRES_PASSWORD=greenops
POSTGRES_DB=greenops
DATABASE_URL=postgres://greenops:greenops@postgres:5432/greenops?sslmode=disable
REDIS_ADDR=redis:6379
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=greenops
```

Proteger le fichier :

```bash
chmod 600 .env
```

Annotation : `.env` contient des secrets. Il ne doit pas etre pousse publiquement.

## 13. Construction des binaires Go

Verifier Go :

```bash
go version
```

Le projet demande Go 1.22 ou plus recent. Si Go n'est pas installe, ou si la version est inferieure a 1.22, installer une version compatible depuis le site officiel Go :

```bash
GO_VERSION="1.22.12"
curl -LO "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
echo 'export PATH=$PATH:/usr/local/go/bin' >> "$HOME/.profile"
source "$HOME/.profile"
go version
```

Compiler les services pour Linux :

```bash
./scripts/build-linux.sh
```

Verifier les binaires :

```bash
ls -lh services/bin
```

Annotation : le script compile `gateway`, `auth-service`, `energy-service` et `alerts-service`.

## 14. Lancement avec Docker Compose

Construire et lancer toute la plateforme :

```bash
docker compose up -d --build
```

Voir les conteneurs :

```bash
docker compose ps
docker ps
```

Voir les images :

```bash
docker image ls
```

Voir les volumes :

```bash
docker volume ls
```

Voir les reseaux :

```bash
docker network ls
```

Annotation : cette etape lance Caddy, les services Go, PostgreSQL, Redis, Prometheus et Grafana.

## 15. Verification Docker Compose

Tester la page principale :

```bash
curl -i http://localhost:8080
```

Tester l'API Gateway :

```bash
curl -s http://localhost:8080/api/health | jq
```

Tester le service energy via le gateway :

```bash
curl -s http://localhost:8080/api/energy/summary | jq
```

Tester le service alerts via le gateway :

```bash
curl -s http://localhost:8080/api/alerts | jq
```

Tester Prometheus :

```bash
curl -i http://localhost:9090/-/ready
```

Tester Grafana :

```bash
curl -s http://localhost:3001/api/health | jq
```

Annotation : si ces commandes repondent, la stack Docker fonctionne.

## 16. Acces depuis le navigateur

Sur la machine physique :

```text
http://localhost:8080
http://localhost:3001
http://localhost:9090
```

Depuis un autre ordinateur du meme reseau :

```bash
hostname -I
```

Noter l'adresse IP principale, puis ouvrir :

```text
http://IP_DE_LA_MACHINE:8080
http://IP_DE_LA_MACHINE:3001
http://IP_DE_LA_MACHINE:9090
```

Compte Grafana de demo :

```text
utilisateur : admin
mot de passe : greenops
```

Annotation : les autres machines ne doivent pas utiliser `localhost`, car `localhost` designe leur propre ordinateur.

## 17. Communication entre conteneurs

Afficher les reseaux Compose :

```bash
docker compose ps
docker network ls
docker network inspect greenops_backend
docker network inspect greenops_data
docker network inspect greenops_monitoring
```

Tester le DNS Docker depuis un conteneur de diagnostic temporaire :

```bash
docker run --rm --network greenops_backend curlimages/curl:latest -s http://auth-service:8080/health
docker run --rm --network greenops_backend curlimages/curl:latest -s http://energy-service:8080/health
docker run --rm --network greenops_backend curlimages/curl:latest -s http://alerts-service:8080/health
```

Tester Redis depuis le conteneur Redis :

```bash
docker compose exec redis redis-cli ping
```

Tester PostgreSQL :

```bash
docker compose exec postgres pg_isready -U greenops -d greenops
```

Annotation : ces commandes prouvent que les conteneurs communiquent par nom de service dans les reseaux Docker.

## 18. Logs Docker utiles

Logs globaux :

```bash
docker compose logs -f
```

Logs d'un service :

```bash
docker compose logs -f gateway
docker compose logs -f auth-service
docker compose logs -f energy-service
docker compose logs -f alerts-service
docker compose logs -f postgres
docker compose logs -f redis
docker compose logs -f prometheus
docker compose logs -f grafana
```

Afficher les derniers logs :

```bash
docker compose logs --tail=100 gateway
```

Annotation : les logs permettent de diagnostiquer un service qui ne demarre pas.

## 19. Arret Docker Compose

Arreter sans supprimer les donnees :

```bash
docker compose stop
```

Redemarrer :

```bash
docker compose start
```

Arreter et supprimer les conteneurs :

```bash
docker compose down
```

Arreter et supprimer aussi les volumes :

```bash
docker compose down -v
```

Annotation : `down -v` supprime les donnees PostgreSQL, Redis, Prometheus et Grafana. A utiliser seulement si on veut repartir de zero.

## 20. Installation Kubernetes avec k3s

Arreter Docker Compose avant Kubernetes si la machine a peu de ressources :

```bash
docker compose down
```

Installer k3s :

```bash
curl -sfL https://get.k3s.io | sh -
```

Verifier le service k3s :

```bash
sudo systemctl status k3s --no-pager
```

Verifier le cluster :

```bash
sudo k3s kubectl get nodes
sudo k3s kubectl get pods -A
```

Installer la configuration Kubernetes pour l'utilisateur courant :

```bash
mkdir -p "$HOME/.kube"
sudo cp /etc/rancher/k3s/k3s.yaml "$HOME/.kube/config"
sudo chown "$USER:$USER" "$HOME/.kube/config"
chmod 600 "$HOME/.kube/config"
```

Verifier `kubectl` :

```bash
kubectl get nodes
kubectl get pods -A
```

Annotation : k3s installe un cluster Kubernetes local mono-noeud, suffisant pour la demo.

## 21. Construction des images pour Kubernetes

Depuis le dossier du projet :

```bash
cd "$HOME/projets/greenops-platform"
```

Construire les images applicatives :

```bash
./scripts/k8s-build-images.sh
```

Verifier les images Docker :

```bash
docker image ls "greenops-*"
```

Annotation : les manifests Kubernetes utilisent les images `greenops-gateway`, `greenops-auth-service`, `greenops-energy-service` et `greenops-alerts-service`.

## 22. Import des images Docker dans k3s

k3s utilise containerd. Les images construites avec Docker doivent donc etre importees dans k3s.

Creer une archive des images :

```bash
docker save -o /tmp/greenops-k8s-images.tar \
  greenops-gateway:latest \
  greenops-auth-service:latest \
  greenops-energy-service:latest \
  greenops-alerts-service:latest
```

Importer dans k3s :

```bash
sudo k3s ctr images import /tmp/greenops-k8s-images.tar
```

Verifier :

```bash
sudo k3s ctr images ls | grep greenops
```

Annotation : sans cette etape, les pods applicatifs peuvent rester en `ImagePullBackOff`.

## 23. Secrets Kubernetes

Pour une demo locale, definir des valeurs simples :

```bash
export JWT_SECRET="change-me-in-production"
export POSTGRES_PASSWORD="greenops"
export GRAFANA_ADMIN_PASSWORD="greenops"
```

Pour un vrai environnement, utiliser des secrets plus longs :

```bash
export JWT_SECRET="remplacer-par-une-valeur-longue"
export POSTGRES_PASSWORD="remplacer-par-un-mot-de-passe-fort"
export GRAFANA_ADMIN_PASSWORD="remplacer-par-un-mot-de-passe-fort"
```

Annotation : le script `k8s-deploy.sh` cree le Secret Kubernetes `greenops-secrets` a partir de ces variables.

## 24. Deploiement Kubernetes

Appliquer le deploiement :

```bash
./scripts/k8s-deploy.sh
```

Lancer le smoke test :

```bash
./scripts/k8s-smoke.sh
```

Voir les objets Kubernetes :

```bash
kubectl get all -n greenops
kubectl get deploy,svc,ingress,hpa,pvc -n greenops
kubectl get pods -n greenops -o wide
```

Annotation : cette etape cree le namespace, les secrets, les deployments, les services, l'ingress, les HPA, les PVC et les network policies.

## 25. Verification Kubernetes

Verifier les rollouts :

```bash
kubectl rollout status deployment/postgres -n greenops
kubectl rollout status deployment/redis -n greenops
kubectl rollout status deployment/auth-service -n greenops
kubectl rollout status deployment/energy-service -n greenops
kubectl rollout status deployment/alerts-service -n greenops
kubectl rollout status deployment/gateway -n greenops
kubectl rollout status deployment/reverse-proxy -n greenops
kubectl rollout status deployment/prometheus -n greenops
kubectl rollout status deployment/grafana -n greenops
```

Voir les logs :

```bash
kubectl logs -n greenops deploy/gateway
kubectl logs -n greenops deploy/auth-service
kubectl logs -n greenops deploy/energy-service
kubectl logs -n greenops deploy/alerts-service
```

Decrire un pod en erreur :

```bash
kubectl describe pod -n greenops NOM_DU_POD
```

Annotation : les rollouts prouvent que Kubernetes a cree les pods attendus.

## 26. Acces Kubernetes par port-forward

Pour tester depuis la machine physique :

```bash
kubectl port-forward -n greenops svc/reverse-proxy 8080:80
```

Dans un autre terminal :

```bash
curl -i http://localhost:8080
```

Pour Grafana :

```bash
kubectl port-forward -n greenops svc/grafana 3001:3000
```

Pour Prometheus :

```bash
kubectl port-forward -n greenops svc/prometheus 9090:9090
```

Annotation : le port-forward est simple pour tester, mais il s'arrete quand le terminal est ferme.

## 27. Exposer Kubernetes sur le reseau local avec port-forward systemd

Creer un service systemd pour l'application :

```bash
sudo tee /etc/systemd/system/greenops-app-portforward.service <<EOF
[Unit]
Description=GreenOps app Kubernetes port-forward
After=k3s.service
Requires=k3s.service

[Service]
User=$USER
Environment=KUBECONFIG=/home/$USER/.kube/config
ExecStart=/usr/local/bin/kubectl port-forward -n greenops svc/reverse-proxy 8080:80 --address 0.0.0.0
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Creer un service pour Grafana :

```bash
sudo tee /etc/systemd/system/greenops-grafana-portforward.service <<EOF
[Unit]
Description=GreenOps Grafana Kubernetes port-forward
After=k3s.service
Requires=k3s.service

[Service]
User=$USER
Environment=KUBECONFIG=/home/$USER/.kube/config
ExecStart=/usr/local/bin/kubectl port-forward -n greenops svc/grafana 3001:3000 --address 0.0.0.0
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Creer un service pour Prometheus :

```bash
sudo tee /etc/systemd/system/greenops-prometheus-portforward.service <<EOF
[Unit]
Description=GreenOps Prometheus Kubernetes port-forward
After=k3s.service
Requires=k3s.service

[Service]
User=$USER
Environment=KUBECONFIG=/home/$USER/.kube/config
ExecStart=/usr/local/bin/kubectl port-forward -n greenops svc/prometheus 9090:9090 --address 0.0.0.0
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Activer les services :

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now greenops-app-portforward.service
sudo systemctl enable --now greenops-grafana-portforward.service
sudo systemctl enable --now greenops-prometheus-portforward.service
```

Verifier :

```bash
sudo systemctl status greenops-app-portforward.service --no-pager
sudo systemctl status greenops-grafana-portforward.service --no-pager
sudo systemctl status greenops-prometheus-portforward.service --no-pager
```

Tester depuis la machine physique :

```bash
curl -i http://localhost:8080
curl -s http://localhost:3001/api/health | jq
curl -i http://localhost:9090/-/ready
```

Tester depuis un autre ordinateur :

```text
http://IP_DE_LA_MACHINE:8080
http://IP_DE_LA_MACHINE:3001
http://IP_DE_LA_MACHINE:9090
```

Annotation : cette methode ressemble a celle utilisee sur la VM pour mettre GreenOps en ligne sur le reseau local.

## 28. Exposer Kubernetes avec Ingress k3s

k3s installe generalement Traefik par defaut. L'ingress du projet utilise ces noms :

```text
greenops.local
prometheus.greenops.local
grafana.greenops.local
```

Sur la machine physique, ajouter ces noms localement :

```bash
sudo tee -a /etc/hosts <<EOF
127.0.0.1 greenops.local
127.0.0.1 prometheus.greenops.local
127.0.0.1 grafana.greenops.local
EOF
```

Depuis un autre ordinateur du meme reseau, ajouter dans son fichier hosts :

```text
IP_DE_LA_MACHINE greenops.local
IP_DE_LA_MACHINE prometheus.greenops.local
IP_DE_LA_MACHINE grafana.greenops.local
```

Tester :

```bash
curl -H "Host: greenops.local" http://IP_DE_LA_MACHINE/
curl -H "Host: prometheus.greenops.local" http://IP_DE_LA_MACHINE/-/ready
curl -H "Host: grafana.greenops.local" http://IP_DE_LA_MACHINE/api/health
```

Annotation : l'ingress est plus proche d'un vrai deploiement, car il utilise des noms DNS.

## 29. Tests applicatifs utiles

Sante globale :

```bash
curl -s http://localhost:8080/api/health | jq
```

Resume energie :

```bash
curl -s http://localhost:8080/api/energy/summary | jq
```

Liste des metriques energie :

```bash
curl -s http://localhost:8080/api/energy/metrics | jq
```

Liste des alertes :

```bash
curl -s http://localhost:8080/api/alerts | jq
```

Evaluation des alertes :

```bash
curl -s http://localhost:8080/api/alerts/evaluate | jq
```

Prometheus pret :

```bash
curl -i http://localhost:9090/-/ready
```

Grafana pret :

```bash
curl -s http://localhost:3001/api/health | jq
```

Annotation : ces commandes peuvent etre montrees pendant une demonstration.

## 30. Test de resilience Kubernetes

Voir les pods gateway :

```bash
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Supprimer un pod gateway :

```bash
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
```

Observer la recreation :

```bash
kubectl get pods -n greenops -w
```

Annotation : Kubernetes recree automatiquement les pods manquants grace au Deployment.

## 31. Test de scaling Kubernetes

Augmenter le nombre de replicas :

```bash
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Revenir a 2 replicas :

```bash
kubectl scale deployment/gateway -n greenops --replicas=2
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Voir les HPA :

```bash
kubectl get hpa -n greenops
```

Annotation : montre la capacite de Kubernetes a augmenter ou reduire le nombre de pods.

## 32. Surveillance disque et ressources

Voir le disque :

```bash
df -h /
```

Voir l'espace Docker :

```bash
docker system df
```

Voir les gros dossiers Docker et k3s :

```bash
sudo du -h -d1 /var/lib/docker | sort -h
sudo du -h -d1 /var/lib/rancher | sort -h
```

Voir CPU et RAM :

```bash
free -h
top
```

Voir les ressources Kubernetes :

```bash
kubectl top nodes
kubectl top pods -n greenops
```

Annotation : `kubectl top` depend du metrics-server. Sur k3s, il est souvent disponible apres quelques minutes.

## 33. Nettoyage Docker

Arreter la stack Docker :

```bash
docker compose down
```

Supprimer les volumes de la stack :

```bash
docker compose down -v
```

Nettoyer les images inutilisees :

```bash
docker image prune -a
```

Nettoyer tout ce qui est inutilise :

```bash
docker system prune -a
```

Annotation : ne pas faire `prune -a` si d'autres projets Docker importants utilisent la machine.

## 34. Nettoyage Kubernetes

Supprimer le namespace GreenOps :

```bash
kubectl delete namespace greenops
```

Arreter les port-forwards systemd :

```bash
sudo systemctl stop greenops-app-portforward.service
sudo systemctl stop greenops-grafana-portforward.service
sudo systemctl stop greenops-prometheus-portforward.service
```

Desactiver les port-forwards :

```bash
sudo systemctl disable greenops-app-portforward.service
sudo systemctl disable greenops-grafana-portforward.service
sudo systemctl disable greenops-prometheus-portforward.service
```

Supprimer les fichiers systemd :

```bash
sudo rm -f /etc/systemd/system/greenops-app-portforward.service
sudo rm -f /etc/systemd/system/greenops-grafana-portforward.service
sudo rm -f /etc/systemd/system/greenops-prometheus-portforward.service
sudo systemctl daemon-reload
```

Arreter k3s :

```bash
sudo systemctl stop k3s
```

Desinstaller k3s si besoin :

```bash
sudo /usr/local/bin/k3s-uninstall.sh
```

Annotation : supprimer k3s efface le cluster local. A utiliser seulement si on veut repartir de zero.

## 35. Diagnostic rapide

Docker ne demarre pas :

```bash
sudo systemctl status docker --no-pager
sudo journalctl -u docker -n 100 --no-pager
```

Permission Docker refusee :

```bash
groups
sudo usermod -aG docker "$USER"
newgrp docker
```

Port deja utilise :

```bash
sudo ss -ltnp | grep -E ':80|:8080|:3001|:9090'
```

Conteneur en erreur :

```bash
docker compose ps
docker compose logs --tail=100 NOM_DU_SERVICE
```

Pod en erreur :

```bash
kubectl get pods -n greenops
kubectl describe pod -n greenops NOM_DU_POD
kubectl logs -n greenops NOM_DU_POD
```

Image Kubernetes non trouvee :

```bash
kubectl describe pod -n greenops NOM_DU_POD
docker image ls "greenops-*"
sudo k3s ctr images ls | grep greenops
```

Puis refaire :

```bash
docker save -o /tmp/greenops-k8s-images.tar \
  greenops-gateway:latest \
  greenops-auth-service:latest \
  greenops-energy-service:latest \
  greenops-alerts-service:latest

sudo k3s ctr images import /tmp/greenops-k8s-images.tar
kubectl rollout restart deployment/gateway -n greenops
kubectl rollout restart deployment/auth-service -n greenops
kubectl rollout restart deployment/energy-service -n greenops
kubectl rollout restart deployment/alerts-service -n greenops
```

Autre machine n'accede pas a l'application :

```bash
ip -br addr
sudo ufw status verbose
sudo ss -ltnp | grep -E ':80|:8080|:3001|:9090'
```

Verifier depuis l'autre machine :

```bash
ping IP_DE_LA_MACHINE
curl -i http://IP_DE_LA_MACHINE:8080
```

Annotation : si `ping` ou `curl` echoue, le probleme vient souvent du reseau, du pare-feu ou d'une mauvaise IP.

## 36. Commandes de demonstration finale

Etat Docker :

```bash
docker compose ps
docker network ls
docker volume ls
```

Etat Kubernetes :

```bash
kubectl get nodes
kubectl get all -n greenops
kubectl get ingress,hpa,pvc -n greenops
```

Tests HTTP :

```bash
curl -s http://localhost:8080/api/health | jq
curl -s http://localhost:8080/api/energy/summary | jq
curl -s http://localhost:8080/api/alerts | jq
curl -s http://localhost:3001/api/health | jq
curl -i http://localhost:9090/-/ready
```

URLs a presenter :

```text
Application Docker ou port-forward Kubernetes : http://IP_DE_LA_MACHINE:8080
Grafana : http://IP_DE_LA_MACHINE:3001
Prometheus : http://IP_DE_LA_MACHINE:9090
Ingress Kubernetes : http://greenops.local
```

## 37. Ordre conseille pour refaire le projet sans se perdre

Pour une demonstration simple :

```bash
sudo apt update
sudo apt install -y git curl jq
git clone https://github.com/montchonabil-collab/greenops-platform.git
cd greenops-platform
cp .env.example .env
chmod +x scripts/*.sh
docker compose config
./scripts/build-linux.sh
docker compose up -d --build
docker compose ps
curl -s http://localhost:8080/api/health | jq
```

Pour la phase Kubernetes :

```bash
docker compose down
curl -sfL https://get.k3s.io | sh -
mkdir -p "$HOME/.kube"
sudo cp /etc/rancher/k3s/k3s.yaml "$HOME/.kube/config"
sudo chown "$USER:$USER" "$HOME/.kube/config"
chmod 600 "$HOME/.kube/config"
./scripts/k8s-build-images.sh
docker save -o /tmp/greenops-k8s-images.tar greenops-gateway:latest greenops-auth-service:latest greenops-energy-service:latest greenops-alerts-service:latest
sudo k3s ctr images import /tmp/greenops-k8s-images.tar
./scripts/k8s-deploy.sh
./scripts/k8s-smoke.sh
kubectl get all -n greenops
```

## 38. Sources officielles consultees

- Docker Engine Ubuntu : https://docs.docker.com/engine/install/ubuntu/
- K3s Quick Start : https://docs.k3s.io/quick-start
- Ubuntu Git : https://documentation.ubuntu.com/ubuntu-for-developers/explanation/use-vcs/
- Go installation : https://go.dev/doc/install/

Ces sources ont ete verifiees le 2026-06-04 pour la partie installation des outils.
