# Installation complete sur une machine physique Windows

Ce document explique comment refaire GreenOps Platform de A a Z sur une machine physique Windows.

Il est adapte a Windows 10 ou Windows 11 avec Docker Desktop, WSL 2, Git, Go, Docker Compose et Kubernetes local via Docker Desktop.

Les commandes sont donnees principalement en PowerShell. Quand une commande doit etre lancee en administrateur, c'est indique clairement.

Important : ne pas mettre de vrais mots de passe dans une documentation partagee. Les valeurs comme `greenops` sont des valeurs de demo locale.

## 1. Objectif final

A la fin, la machine Windows doit pouvoir lancer :

- GreenOps avec Docker Compose sur `http://localhost:8080` ;
- Grafana sur `http://localhost:3001` ;
- Prometheus sur `http://localhost:9090` ;
- GreenOps dans Kubernetes local avec Docker Desktop ;
- l'acces aux services depuis une autre machine du meme reseau ;
- des commandes de verification pour prouver que les conteneurs et les pods communiquent.

## 2. Architecture cible sur Windows

Sur une machine Windows, le plus simple est d'utiliser :

- Windows comme poste physique principal ;
- WSL 2 comme backend Linux de Docker Desktop ;
- Docker Desktop pour Docker Engine, Docker Compose et Kubernetes ;
- Git for Windows pour recuperer le projet ;
- Go pour compiler les services en binaires Linux ;
- PowerShell pour executer les commandes ;
- Docker Compose pour la premiere phase ;
- Kubernetes Docker Desktop pour la phase orchestration.

## 3. Prerequis materiels et systeme

Machine conseillee :

- Windows 10 22H2 ou Windows 11 ;
- processeur 64 bits ;
- virtualisation activee dans le BIOS ou UEFI ;
- 8 Go RAM minimum conseille ;
- 30 Go disque libre minimum conseille ;
- acces Internet ;
- droits administrateur pour installer les outils ;
- meme reseau local si d'autres utilisateurs doivent acceder a l'application.

Ports utilises :

- `8080` pour l'application ;
- `3001` pour Grafana ;
- `9090` pour Prometheus ;
- `80` si on utilise un Ingress HTTP.

## 4. Ouvrir PowerShell

Pour les commandes normales :

```powershell
PowerShell
```

Pour les commandes administrateur :

```text
Clic droit sur Windows Terminal ou PowerShell
Executer en tant qu'administrateur
```

Annotation : l'installation de WSL, Docker Desktop et les regles pare-feu demandent souvent les droits administrateur.

## 5. Verifier Windows

Afficher la version Windows :

```powershell
winver
```

Afficher les informations systeme :

```powershell
systeminfo
```

Verifier la version PowerShell :

```powershell
$PSVersionTable
```

Verifier le disque :

```powershell
Get-PSDrive -PSProvider FileSystem | Select-Object Name,Used,Free,Root
```

Verifier la RAM et le CPU :

```powershell
Get-CimInstance Win32_ComputerSystem | Select-Object Manufacturer,Model,TotalPhysicalMemory
Get-CimInstance Win32_Processor | Select-Object Name,NumberOfCores,NumberOfLogicalProcessors
```

Annotation : Docker Desktop et Kubernetes local ont besoin de ressources suffisantes.

## 6. Verifier la virtualisation

Afficher les informations Hyper-V :

```powershell
systeminfo | Select-String -Pattern "Hyper-V"
```

Verifier les fonctions Windows utiles :

```powershell
Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux
Get-WindowsOptionalFeature -Online -FeatureName VirtualMachinePlatform
```

Annotation : WSL 2 et Docker Desktop ont besoin de la virtualisation activee.

Si la virtualisation n'est pas activee, il faut l'activer dans le BIOS ou UEFI de la machine physique.

## 7. Installer ou verifier winget

Verifier `winget` :

```powershell
winget --version
```

Mettre a jour les sources :

```powershell
winget source update
```

Annotation : `winget` sert a installer Git, Go, Docker Desktop et GitHub CLI rapidement.

Si `winget` n'existe pas, installer ou mettre a jour "App Installer" depuis le Microsoft Store.

## 8. Installer WSL 2

Dans PowerShell administrateur :

```powershell
wsl --install
```

Redemarrer Windows :

```powershell
Restart-Computer
```

Apres redemarrage, verifier WSL :

```powershell
wsl --status
wsl --list --verbose
```

Forcer WSL 2 comme version par defaut :

```powershell
wsl --set-default-version 2
```

Mettre a jour WSL :

```powershell
wsl --update
```

Annotation : Docker Desktop utilise WSL 2 pour executer les conteneurs Linux sur Windows.

## 9. Installer Docker Desktop

Methode simple avec winget :

```powershell
winget install --id Docker.DockerDesktop -e --source winget
```

Alternative avec l'installeur officiel :

```powershell
Start-Process ".\Docker Desktop Installer.exe" -Wait -ArgumentList "install --accept-license"
```

Demarrer Docker Desktop :

```powershell
Start-Process "Docker Desktop"
```

Si la commande precedente ne trouve pas l'application, essayer :

```powershell
Start-Process "$Env:ProgramFiles\Docker\Docker\Docker Desktop.exe"
```

Ou en installation par utilisateur :

```powershell
Start-Process "$Env:LOCALAPPDATA\Programs\DockerDesktop\Docker Desktop.exe"
```

Annotation : Docker Desktop doit etre ouvert pour que `docker` et `docker compose` fonctionnent.

## 10. Verifier Docker Desktop

Ouvrir un nouveau PowerShell, puis verifier :

```powershell
docker version
docker compose version
docker info
```

Tester avec une image simple :

```powershell
docker run --rm hello-world
```

Voir les contextes Docker :

```powershell
docker context ls
```

Annotation : le contexte actif doit normalement pointer vers Docker Desktop.

## 11. Ajouter l'utilisateur au groupe docker-users

Si Docker Desktop a ete installe pour tous les utilisateurs, il peut etre necessaire d'ajouter le compte Windows au groupe `docker-users`.

Dans PowerShell administrateur :

```powershell
net localgroup docker-users "$env:USERNAME" /add
```

Se deconnecter puis se reconnecter.

Verifier :

```powershell
whoami /groups
```

Annotation : cette etape evite certaines erreurs de permission avec Docker Desktop.

## 12. Installer Git for Windows

Installer Git :

```powershell
winget install --id Git.Git -e --source winget
```

Fermer puis rouvrir PowerShell.

Verifier Git :

```powershell
git --version
where.exe git
```

Configurer Git :

```powershell
git config --global user.name "NOM_GIT"
git config --global user.email "EMAIL_GIT"
git config --global init.defaultBranch main
```

Verifier :

```powershell
git config --global --list
```

Annotation : Git sert a recuperer le code, committer et pousser vers GitHub.

## 13. Installer GitHub CLI

Installer GitHub CLI :

```powershell
winget install --id GitHub.cli -e --source winget
```

Fermer puis rouvrir PowerShell.

Verifier :

```powershell
gh --version
```

Se connecter a GitHub si besoin :

```powershell
gh auth login
```

Verifier la connexion :

```powershell
gh auth status
gh api user
```

Annotation : GitHub CLI est utile pour authentifier Git et consulter les workflows CI.

## 14. Installer Go

Installer Go avec winget :

```powershell
winget install --id GoLang.Go -e --source winget
```

Fermer puis rouvrir PowerShell.

Verifier Go :

```powershell
go version
where.exe go
```

Le projet utilise Go 1.22 ou plus recent. Verifier que la version affichee est compatible.

Annotation : Go sert a compiler les microservices avant de construire les images Docker.

## 15. Installer Visual Studio Code optionnel

Installer VS Code :

```powershell
winget install --id Microsoft.VisualStudioCode -e --source winget
```

Ouvrir le projet plus tard :

```powershell
code C:\docker-projets\greenops-platform
```

Annotation : VS Code est optionnel mais pratique pour lire et modifier les fichiers.

## 16. Creer le dossier de travail

Creer un dossier simple sur le disque C :

```powershell
New-Item -ItemType Directory -Path C:\docker-projets -Force
Set-Location C:\docker-projets
```

Verifier :

```powershell
Get-Location
Get-ChildItem
```

Annotation : `C:\docker-projets` evite les chemins trop longs et garde les projets Docker au meme endroit.

## 17. Recuperer le projet depuis GitHub

Cloner le depot :

```powershell
Set-Location C:\docker-projets
git clone https://github.com/montchonabil-collab/greenops-platform.git
Set-Location C:\docker-projets\greenops-platform
```

Verifier :

```powershell
git status
git log --oneline --decorate -5
```

Annotation : le projet est maintenant sur la machine Windows.

## 18. Alternative sans git clone

Si GitHub n'est pas accessible avec `git clone`, telecharger l'archive :

```powershell
Set-Location C:\docker-projets
Invoke-WebRequest -Uri "https://github.com/montchonabil-collab/greenops-platform/archive/refs/heads/main.zip" -OutFile "greenops-platform.zip"
Expand-Archive -Path "greenops-platform.zip" -DestinationPath .
Rename-Item -Path "greenops-platform-main" -NewName "greenops-platform"
Set-Location C:\docker-projets\greenops-platform
```

Verifier :

```powershell
Get-ChildItem
```

Annotation : cette methode permet de recuperer le projet sans configuration Git.

## 19. Verifier la structure du projet

Afficher les fichiers principaux :

```powershell
Get-ChildItem
Get-ChildItem -Recurse -Depth 2 | Select-Object FullName
```

Voir les dossiers importants :

```powershell
Get-ChildItem .\services
Get-ChildItem .\frontend
Get-ChildItem .\docker
Get-ChildItem .\k8s
Get-ChildItem .\monitoring
```

Annotation : on doit voir `compose.yaml`, `services`, `frontend`, `k8s`, `monitoring` et `docs`.

## 20. Configurer le fichier .env

Copier le fichier exemple :

```powershell
Copy-Item .env.example .env
```

Ouvrir dans Notepad :

```powershell
notepad .env
```

Exemple de contenu pour une demo :

```text
PROJECT_NAME=greenops
JWT_SECRET=GreenOpsJwtDemo2026u7N9pQ2sV5xL8rT4mK6c
POSTGRES_USER=greenops
POSTGRES_PASSWORD=GreenOpsPgDemo2026V7nR4qT9sL2
POSTGRES_DB=greenops
DATABASE_URL=postgres://greenops:GreenOpsPgDemo2026V7nR4qT9sL2@postgres:5432/greenops?sslmode=disable
REDIS_ADDR=redis:6379
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=GreenOpsGrafanaDemo2026Q8mK5pZ1vX
```

Verifier que le fichier existe :

```powershell
Get-Item .env
```

Annotation : `.env` est utilise par Docker Compose pour injecter les variables dans les services.

## 21. Tester le code Go

Aller dans le dossier des services :

```powershell
Set-Location C:\docker-projets\greenops-platform\services
```

Telecharger les dependances :

```powershell
go mod download
```

Lancer les tests :

```powershell
go test ./...
```

Revenir a la racine :

```powershell
Set-Location C:\docker-projets\greenops-platform
```

Annotation : les tests Go permettent de verifier que le code compile correctement avant Docker.

## 22. Compiler les binaires Linux depuis Windows

Les Dockerfiles du projet utilisent des images `scratch`. Ils copient donc des binaires Linux deja compiles dans `services\bin`.

Aller dans les services :

```powershell
Set-Location C:\docker-projets\greenops-platform\services
```

Creer le dossier `bin` :

```powershell
New-Item -ItemType Directory -Path .\bin -Force
```

Configurer la compilation Linux :

```powershell
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Compiler le gateway :

```powershell
go build -trimpath -ldflags="-s -w" -o .\bin\gateway .\gateway
```

Compiler le service auth :

```powershell
go build -trimpath -ldflags="-s -w" -o .\bin\auth-service .\auth
```

Compiler le service energy :

```powershell
go build -trimpath -ldflags="-s -w" -o .\bin\energy-service .\energy
```

Compiler le service alerts :

```powershell
go build -trimpath -ldflags="-s -w" -o .\bin\alerts-service .\alerts
```

Verifier les binaires :

```powershell
Get-ChildItem .\bin
```

Nettoyer les variables de session si besoin :

```powershell
Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
```

Revenir a la racine :

```powershell
Set-Location C:\docker-projets\greenops-platform
```

Annotation : cette etape est indispensable sur Windows, car les conteneurs executent Linux, pas des binaires Windows.

## 23. Valider Docker Compose

Verifier la configuration Compose :

```powershell
docker compose config
```

Annotation : cette commande detecte les erreurs YAML et affiche la configuration finale.

## 24. Construire et lancer Docker Compose

Lancer toute la plateforme :

```powershell
docker compose up -d --build
```

Voir les conteneurs Compose :

```powershell
docker compose ps
```

Voir tous les conteneurs Docker :

```powershell
docker ps
```

Voir les images :

```powershell
docker image ls
```

Voir les volumes :

```powershell
docker volume ls
```

Voir les reseaux :

```powershell
docker network ls
```

Annotation : cette etape lance Caddy, Gateway, Auth, Energy, Alerts, PostgreSQL, Redis, Prometheus et Grafana.

## 25. Tester l'application Docker Compose

Tester la page principale :

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080
```

Tester l'API Gateway :

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/health
```

Tester le resume energie :

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/energy/summary
```

Tester les alertes :

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/alerts
```

Tester Prometheus :

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:9090/-/ready
```

Tester Grafana :

```powershell
Invoke-RestMethod -Uri http://localhost:3001/api/health
```

Annotation : si ces commandes repondent, la stack Docker fonctionne.

## 26. Ouvrir dans le navigateur

Application :

```powershell
Start-Process "http://localhost:8080"
```

Grafana :

```powershell
Start-Process "http://localhost:3001"
```

Prometheus :

```powershell
Start-Process "http://localhost:9090"
```

Compte Grafana de demo :

```text
utilisateur : admin
mot de passe : GreenOpsGrafanaDemo2026Q8mK5pZ1vX
```

Annotation : sur la machine Windows, `localhost` pointe vers le PC lui-meme.

## 27. Voir les logs Docker

Voir tous les logs :

```powershell
docker compose logs -f
```

Voir les logs du gateway :

```powershell
docker compose logs -f gateway
```

Voir les logs du service energy :

```powershell
docker compose logs -f energy-service
```

Voir les logs de PostgreSQL :

```powershell
docker compose logs -f postgres
```

Voir les derniers logs seulement :

```powershell
docker compose logs --tail=100 gateway
```

Annotation : les logs aident a trouver rapidement un service en erreur.

## 28. Verifier la communication entre conteneurs

Afficher les reseaux :

```powershell
docker network ls
```

Inspecter le reseau backend :

```powershell
docker network inspect greenops_backend
```

Tester le DNS Docker avec un conteneur de diagnostic :

```powershell
docker run --rm --network greenops_backend curlimages/curl:latest -s http://auth-service:8080/health
docker run --rm --network greenops_backend curlimages/curl:latest -s http://energy-service:8080/health
docker run --rm --network greenops_backend curlimages/curl:latest -s http://alerts-service:8080/health
```

Tester Redis :

```powershell
docker compose exec redis redis-cli ping
```

Tester PostgreSQL :

```powershell
docker compose exec postgres pg_isready -U greenops -d greenops
```

Annotation : ces commandes prouvent que les conteneurs communiquent par nom de service dans Docker.

## 29. Acces depuis un autre ordinateur du meme reseau

Trouver l'adresse IP Windows :

```powershell
ipconfig
```

Ou avec PowerShell :

```powershell
Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -notlike "127.*"} | Select-Object InterfaceAlias,IPAddress
```

Noter l'IP du reseau principal, par exemple :

```text
IP_DE_LA_MACHINE_WINDOWS
```

Depuis l'autre ordinateur, ouvrir :

```text
http://IP_DE_LA_MACHINE_WINDOWS:8080
http://IP_DE_LA_MACHINE_WINDOWS:3001
http://IP_DE_LA_MACHINE_WINDOWS:9090
```

Annotation : les autres ordinateurs ne doivent pas utiliser `localhost`, car leur `localhost` correspond a leur propre machine.

## 30. Ouvrir le pare-feu Windows

Dans PowerShell administrateur :

```powershell
New-NetFirewallRule -DisplayName "GreenOps App 8080" -Direction Inbound -Protocol TCP -LocalPort 8080 -Action Allow
New-NetFirewallRule -DisplayName "GreenOps Grafana 3001" -Direction Inbound -Protocol TCP -LocalPort 3001 -Action Allow
New-NetFirewallRule -DisplayName "GreenOps Prometheus 9090" -Direction Inbound -Protocol TCP -LocalPort 9090 -Action Allow
```

Verifier :

```powershell
Get-NetFirewallRule -DisplayName "GreenOps*"
```

Tester localement :

```powershell
Test-NetConnection -ComputerName localhost -Port 8080
Test-NetConnection -ComputerName localhost -Port 3001
Test-NetConnection -ComputerName localhost -Port 9090
```

Tester depuis un autre PC :

```powershell
Test-NetConnection -ComputerName IP_DE_LA_MACHINE_WINDOWS -Port 8080
```

Annotation : si un autre ordinateur ne peut pas acceder a GreenOps, verifier l'IP, le pare-feu et le reseau.

## 31. Arreter Docker Compose

Arreter les conteneurs sans supprimer les donnees :

```powershell
docker compose stop
```

Redemarrer :

```powershell
docker compose start
```

Arreter et supprimer les conteneurs :

```powershell
docker compose down
```

Arreter et supprimer aussi les volumes :

```powershell
docker compose down -v
```

Annotation : `down -v` supprime les donnees PostgreSQL, Redis, Prometheus et Grafana.

## 32. Activer Kubernetes dans Docker Desktop

Ouvrir Docker Desktop.

Depuis l'interface :

```text
Docker Desktop
Kubernetes
Create cluster
Choisir Kubeadm ou kind
Create
```

Pour une demo simple, un cluster mono-noeud suffit.

Verifier que `kubectl` est disponible :

```powershell
kubectl version --client
kubectl config get-contexts
```

Utiliser le contexte Docker Desktop :

```powershell
kubectl config use-context docker-desktop
```

Verifier le noeud :

```powershell
kubectl get nodes
kubectl get pods -A
```

Annotation : Docker Desktop installe `kubectl` avec son cluster Kubernetes local.

## 33. Si kubectl n'est pas trouve

Chercher `kubectl` :

```powershell
where.exe kubectl
```

Chemin possible avec installation tous utilisateurs :

```powershell
& "C:\Program Files\Docker\Docker\resources\bin\kubectl.exe" version --client
```

Chemin possible avec installation par utilisateur :

```powershell
& "$env:LOCALAPPDATA\Programs\DockerDesktop\resources\bin\kubectl.exe" version --client
```

Annotation : si `kubectl` existe mais n'est pas dans le PATH, rouvrir PowerShell ou ajouter le dossier Docker Desktop au PATH.

## 34. Preparer les images pour Kubernetes

Arreter Docker Compose pour liberer les ports :

```powershell
docker compose down
```

Verifier que les binaires Linux existent :

```powershell
Get-ChildItem C:\docker-projets\greenops-platform\services\bin
```

Si les binaires n'existent pas, refaire la compilation :

```powershell
Set-Location C:\docker-projets\greenops-platform\services
New-Item -ItemType Directory -Path .\bin -Force
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -trimpath -ldflags="-s -w" -o .\bin\gateway .\gateway
go build -trimpath -ldflags="-s -w" -o .\bin\auth-service .\auth
go build -trimpath -ldflags="-s -w" -o .\bin\energy-service .\energy
go build -trimpath -ldflags="-s -w" -o .\bin\alerts-service .\alerts
Set-Location C:\docker-projets\greenops-platform
```

Construire les images applicatives :

```powershell
docker build -t greenops-gateway:latest -f .\services\gateway\Dockerfile .\services
docker build -t greenops-auth-service:latest -f .\services\auth\Dockerfile .\services
docker build -t greenops-energy-service:latest -f .\services\energy\Dockerfile .\services
docker build -t greenops-alerts-service:latest -f .\services\alerts\Dockerfile .\services
```

Verifier :

```powershell
docker image ls "greenops-*"
```

Annotation : les manifests Kubernetes utilisent ces images locales.

## 35. Deployer les secrets Kubernetes

Definir les valeurs de demo dans PowerShell :

```powershell
$env:JWT_SECRET = "GreenOpsJwtDemo2026u7N9pQ2sV5xL8rT4mK6c"
$env:POSTGRES_PASSWORD = "GreenOpsPgDemo2026V7nR4qT9sL2"
$env:GRAFANA_ADMIN_PASSWORD = "GreenOpsGrafanaDemo2026Q8mK5pZ1vX"
```

Creer le namespace :

```powershell
kubectl apply -f .\k8s\namespace.yaml
```

Creer le secret :

```powershell
kubectl create secret generic greenops-secrets `
  --namespace greenops `
  --from-literal=JWT_SECRET="$env:JWT_SECRET" `
  --from-literal=POSTGRES_PASSWORD="$env:POSTGRES_PASSWORD" `
  --from-literal=GRAFANA_ADMIN_PASSWORD="$env:GRAFANA_ADMIN_PASSWORD" `
  --dry-run=client `
  -o yaml | kubectl apply -f -
```

Verifier :

```powershell
kubectl get secret -n greenops
```

Annotation : le secret injecte les mots de passe dans les pods sans les ecrire dans les manifests.

## 36. Deployer les manifests Kubernetes

Rendre les manifests avec Kustomize :

```powershell
kubectl kustomize --load-restrictor=LoadRestrictionsNone .\k8s
```

Appliquer les manifests :

```powershell
kubectl kustomize --load-restrictor=LoadRestrictionsNone .\k8s | kubectl apply -f -
```

Verifier les objets :

```powershell
kubectl get deploy,svc,ingress,hpa,pvc -n greenops
kubectl get pods -n greenops -o wide
```

Annotation : cette etape cree les deployments, services, ingress, HPA, PVC, configmaps et network policies.

## 37. Attendre les rollouts Kubernetes

Attendre PostgreSQL :

```powershell
kubectl rollout status deployment/postgres -n greenops --timeout=180s
```

Attendre Redis :

```powershell
kubectl rollout status deployment/redis -n greenops --timeout=180s
```

Attendre les services applicatifs :

```powershell
kubectl rollout status deployment/auth-service -n greenops --timeout=180s
kubectl rollout status deployment/energy-service -n greenops --timeout=180s
kubectl rollout status deployment/alerts-service -n greenops --timeout=180s
kubectl rollout status deployment/gateway -n greenops --timeout=180s
kubectl rollout status deployment/reverse-proxy -n greenops --timeout=180s
```

Attendre le monitoring :

```powershell
kubectl rollout status deployment/prometheus -n greenops --timeout=180s
kubectl rollout status deployment/grafana -n greenops --timeout=180s
```

Annotation : les rollouts confirment que Kubernetes a cree les pods attendus.

## 38. Smoke test Kubernetes

Voir les pods :

```powershell
kubectl get pods -n greenops -o wide
```

Voir les services :

```powershell
kubectl get svc -n greenops
```

Attendre les pods principaux :

```powershell
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=gateway -n greenops --timeout=120s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=reverse-proxy -n greenops --timeout=120s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=grafana -n greenops --timeout=120s
```

Annotation : c'est l'equivalent PowerShell du script `scripts/k8s-smoke.sh`.

## 39. Acceder a Kubernetes avec port-forward

Ouvrir un terminal PowerShell pour l'application :

```powershell
kubectl port-forward -n greenops svc/reverse-proxy 8080:80 --address 0.0.0.0
```

Ouvrir un deuxieme terminal pour Grafana :

```powershell
kubectl port-forward -n greenops svc/grafana 3001:3000 --address 0.0.0.0
```

Ouvrir un troisieme terminal pour Prometheus :

```powershell
kubectl port-forward -n greenops svc/prometheus 9090:9090 --address 0.0.0.0
```

Tester dans un autre PowerShell :

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080
Invoke-RestMethod -Uri http://localhost:3001/api/health
Invoke-WebRequest -UseBasicParsing http://localhost:9090/-/ready
```

Annotation : `--address 0.0.0.0` permet aussi l'acces depuis une autre machine du reseau, si le pare-feu Windows autorise les ports.

## 40. Lancer les port-forwards dans des fenetres separees

Application :

```powershell
Start-Process powershell -ArgumentList "-NoExit", "-Command", "kubectl port-forward -n greenops svc/reverse-proxy 8080:80 --address 0.0.0.0"
```

Grafana :

```powershell
Start-Process powershell -ArgumentList "-NoExit", "-Command", "kubectl port-forward -n greenops svc/grafana 3001:3000 --address 0.0.0.0"
```

Prometheus :

```powershell
Start-Process powershell -ArgumentList "-NoExit", "-Command", "kubectl port-forward -n greenops svc/prometheus 9090:9090 --address 0.0.0.0"
```

Annotation : pratique pour une demo, car chaque port-forward reste visible dans sa fenetre.

## 41. Acces Kubernetes depuis une autre machine

Trouver l'IP Windows :

```powershell
Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -notlike "127.*"} | Select-Object InterfaceAlias,IPAddress
```

Ouvrir les ports dans le pare-feu si ce n'est pas deja fait :

```powershell
New-NetFirewallRule -DisplayName "GreenOps K8s App 8080" -Direction Inbound -Protocol TCP -LocalPort 8080 -Action Allow
New-NetFirewallRule -DisplayName "GreenOps K8s Grafana 3001" -Direction Inbound -Protocol TCP -LocalPort 3001 -Action Allow
New-NetFirewallRule -DisplayName "GreenOps K8s Prometheus 9090" -Direction Inbound -Protocol TCP -LocalPort 9090 -Action Allow
```

Depuis un autre ordinateur :

```text
http://IP_DE_LA_MACHINE_WINDOWS:8080
http://IP_DE_LA_MACHINE_WINDOWS:3001
http://IP_DE_LA_MACHINE_WINDOWS:9090
```

Annotation : cette methode correspond a la mise en ligne locale du projet pour les autres utilisateurs de la salle.

## 42. Option Ingress Kubernetes

Le projet contient deja un ingress avec ces noms :

```text
greenops.local
prometheus.greenops.local
grafana.greenops.local
```

Afficher l'ingress :

```powershell
kubectl get ingress -n greenops
kubectl describe ingress greenops -n greenops
```

Ajouter les noms dans le fichier hosts Windows.

Ouvrir Notepad en administrateur :

```powershell
Start-Process notepad "C:\Windows\System32\drivers\etc\hosts" -Verb RunAs
```

Ajouter :

```text
127.0.0.1 greenops.local
127.0.0.1 prometheus.greenops.local
127.0.0.1 grafana.greenops.local
```

Tester :

```powershell
Invoke-WebRequest -UseBasicParsing http://greenops.local
Invoke-WebRequest -UseBasicParsing http://prometheus.greenops.local/-/ready
Invoke-RestMethod -Uri http://grafana.greenops.local/api/health
```

Annotation : si aucun Ingress Controller n'est disponible dans Docker Desktop, utiliser les port-forwards de la section precedente.

## 43. Verifier la communication entre pods

Voir les services internes :

```powershell
kubectl get svc -n greenops
```

Lancer un pod de diagnostic temporaire :

```powershell
kubectl run curl-test -n greenops --rm -it --image=curlimages/curl:latest --restart=Never -- http://gateway:8080/api/health
```

Tester Auth :

```powershell
kubectl run curl-auth -n greenops --rm -it --image=curlimages/curl:latest --restart=Never -- http://auth-service:8080/health
```

Tester Energy :

```powershell
kubectl run curl-energy -n greenops --rm -it --image=curlimages/curl:latest --restart=Never -- http://energy-service:8080/health
```

Tester Alerts :

```powershell
kubectl run curl-alerts -n greenops --rm -it --image=curlimages/curl:latest --restart=Never -- http://alerts-service:8080/health
```

Annotation : ces commandes prouvent que Kubernetes resout les noms de Services et que les pods communiquent.

## 44. Voir les logs Kubernetes

Gateway :

```powershell
kubectl logs -n greenops deploy/gateway
```

Auth :

```powershell
kubectl logs -n greenops deploy/auth-service
```

Energy :

```powershell
kubectl logs -n greenops deploy/energy-service
```

Alerts :

```powershell
kubectl logs -n greenops deploy/alerts-service
```

Grafana :

```powershell
kubectl logs -n greenops deploy/grafana
```

Prometheus :

```powershell
kubectl logs -n greenops deploy/prometheus
```

Annotation : utile si un pod est en `CrashLoopBackOff` ou `ImagePullBackOff`.

## 45. Diagnostiquer un pod en erreur

Voir les pods :

```powershell
kubectl get pods -n greenops
```

Decrire un pod :

```powershell
kubectl describe pod -n greenops NOM_DU_POD
```

Voir les evenements :

```powershell
kubectl get events -n greenops --sort-by=.lastTimestamp
```

Voir les images locales :

```powershell
docker image ls "greenops-*"
```

Redemarrer les deployments applicatifs :

```powershell
kubectl rollout restart deployment/gateway -n greenops
kubectl rollout restart deployment/auth-service -n greenops
kubectl rollout restart deployment/energy-service -n greenops
kubectl rollout restart deployment/alerts-service -n greenops
```

Annotation : `ImagePullBackOff` signifie souvent que l'image locale n'est pas disponible pour le cluster.

## 46. Test de resilience Kubernetes

Voir les pods gateway :

```powershell
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Supprimer un pod gateway :

```powershell
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
```

Observer la recreation :

```powershell
kubectl get pods -n greenops -w
```

Annotation : Kubernetes recree automatiquement les pods manquants grace au Deployment.

## 47. Test de scaling Kubernetes

Passer le gateway a 4 replicas :

```powershell
kubectl scale deployment/gateway -n greenops --replicas=4
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Revenir a 2 replicas :

```powershell
kubectl scale deployment/gateway -n greenops --replicas=2
kubectl get pods -n greenops -l app.kubernetes.io/name=gateway
```

Voir les HPA :

```powershell
kubectl get hpa -n greenops
```

Annotation : montre que Kubernetes peut augmenter ou reduire le nombre de pods.

## 48. Commandes de demonstration Docker

Etat de la stack :

```powershell
docker compose ps
```

Services publies :

```powershell
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

Reseaux :

```powershell
docker network ls
```

Volumes :

```powershell
docker volume ls
```

Tests HTTP :

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/health
Invoke-RestMethod -Uri http://localhost:8080/api/energy/summary
Invoke-RestMethod -Uri http://localhost:8080/api/alerts
```

## 49. Commandes de demonstration Kubernetes

Etat du cluster :

```powershell
kubectl get nodes
kubectl get all -n greenops
```

Objets importants :

```powershell
kubectl get deploy,svc,ingress,hpa,pvc -n greenops
```

Pods avec leur noeud :

```powershell
kubectl get pods -n greenops -o wide
```

Logs gateway :

```powershell
kubectl logs -n greenops deploy/gateway
```

Resilience :

```powershell
kubectl delete pod -n greenops -l app.kubernetes.io/name=gateway
kubectl get pods -n greenops -w
```

## 50. Nettoyage Docker

Depuis la racine du projet :

```powershell
Set-Location C:\docker-projets\greenops-platform
docker compose down
```

Supprimer les volumes :

```powershell
docker compose down -v
```

Nettoyer les images inutilisees :

```powershell
docker image prune -a
```

Nettoyer tout ce qui est inutilise :

```powershell
docker system prune -a
```

Annotation : attention, `prune -a` peut supprimer des images d'autres projets Docker.

## 51. Nettoyage Kubernetes

Supprimer GreenOps :

```powershell
kubectl delete namespace greenops
```

Verifier :

```powershell
kubectl get namespaces
kubectl get pods -A
```

Annotation : supprimer le namespace efface les pods, services, PVC et autres objets GreenOps dans Kubernetes.

## 52. Supprimer les regles pare-feu GreenOps

Dans PowerShell administrateur :

```powershell
Get-NetFirewallRule -DisplayName "GreenOps*" | Remove-NetFirewallRule
```

Verifier :

```powershell
Get-NetFirewallRule -DisplayName "GreenOps*"
```

Annotation : a faire seulement si on ne veut plus exposer l'application sur le reseau.

## 53. Diagnostic Docker Desktop

Verifier si Docker repond :

```powershell
docker info
```

Voir les processus Docker :

```powershell
Get-Process | Where-Object {$_.ProcessName -like "*Docker*"}
```

Relancer Docker Desktop :

```powershell
Start-Process "Docker Desktop"
```

Arreter WSL si Docker Desktop bloque :

```powershell
wsl --shutdown
```

Puis relancer Docker Desktop :

```powershell
Start-Process "Docker Desktop"
```

Annotation : `wsl --shutdown` redemarre le backend Linux utilise par Docker Desktop.

## 54. Diagnostic reseau Windows

Verifier les ports ouverts :

```powershell
Get-NetTCPConnection -LocalPort 8080,3001,9090 -ErrorAction SilentlyContinue
```

Tester un port local :

```powershell
Test-NetConnection -ComputerName localhost -Port 8080
```

Tester depuis la machine voisine :

```powershell
Test-NetConnection -ComputerName IP_DE_LA_MACHINE_WINDOWS -Port 8080
```

Afficher les IP :

```powershell
ipconfig
```

Annotation : si l'autre machine ne voit pas l'application, verifier l'IP, le pare-feu et que le port-forward ou Docker Compose est bien actif.

## 55. Probleme courant : port deja utilise

Chercher le processus qui utilise le port 8080 :

```powershell
Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue | Select-Object LocalAddress,LocalPort,State,OwningProcess
```

Identifier le processus :

```powershell
Get-Process -Id ID_DU_PROCESSUS
```

Arreter Docker Compose si necessaire :

```powershell
docker compose down
```

Annotation : Docker Compose et `kubectl port-forward` ne peuvent pas utiliser le meme port en meme temps.

## 56. Probleme courant : les images Kubernetes ne partent pas

Voir le statut :

```powershell
kubectl get pods -n greenops
```

Decrire le pod :

```powershell
kubectl describe pod -n greenops NOM_DU_POD
```

Verifier les images Docker :

```powershell
docker image ls "greenops-*"
```

Reconstruire les images :

```powershell
docker build -t greenops-gateway:latest -f .\services\gateway\Dockerfile .\services
docker build -t greenops-auth-service:latest -f .\services\auth\Dockerfile .\services
docker build -t greenops-energy-service:latest -f .\services\energy\Dockerfile .\services
docker build -t greenops-alerts-service:latest -f .\services\alerts\Dockerfile .\services
```

Redemarrer les deployments :

```powershell
kubectl rollout restart deployment/gateway -n greenops
kubectl rollout restart deployment/auth-service -n greenops
kubectl rollout restart deployment/energy-service -n greenops
kubectl rollout restart deployment/alerts-service -n greenops
```

Annotation : les images doivent etre construites dans le meme Docker Desktop que celui utilise par le cluster Kubernetes.

## 57. Probleme courant : Go compile des binaires Windows

Verifier les variables :

```powershell
Get-ChildItem Env:GOOS
Get-ChildItem Env:GOARCH
Get-ChildItem Env:CGO_ENABLED
```

Remettre la compilation Linux :

```powershell
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Recompiler :

```powershell
Set-Location C:\docker-projets\greenops-platform\services
go build -trimpath -ldflags="-s -w" -o .\bin\gateway .\gateway
go build -trimpath -ldflags="-s -w" -o .\bin\auth-service .\auth
go build -trimpath -ldflags="-s -w" -o .\bin\energy-service .\energy
go build -trimpath -ldflags="-s -w" -o .\bin\alerts-service .\alerts
```

Annotation : les conteneurs GreenOps executent Linux. Il faut donc des binaires Linux meme si la compilation se fait depuis Windows.

## 58. Sauvegarder le travail avec Git

Voir les modifications :

```powershell
git status
git diff --stat
```

Ajouter les fichiers :

```powershell
git add .
```

Committer :

```powershell
git commit -m "Update GreenOps project"
```

Pousser :

```powershell
git push
```

Voir les workflows GitHub :

```powershell
gh run list --repo montchonabil-collab/greenops-platform --limit 5
```

Annotation : cette partie sert si on modifie le projet depuis la machine Windows.

## 59. Ordre rapide pour refaire Docker Compose

Resume des commandes principales :

```powershell
Set-Location C:\docker-projets
git clone https://github.com/montchonabil-collab/greenops-platform.git
Set-Location C:\docker-projets\greenops-platform
Copy-Item .env.example .env
Set-Location .\services
New-Item -ItemType Directory -Path .\bin -Force
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go mod download
go test ./...
go build -trimpath -ldflags="-s -w" -o .\bin\gateway .\gateway
go build -trimpath -ldflags="-s -w" -o .\bin\auth-service .\auth
go build -trimpath -ldflags="-s -w" -o .\bin\energy-service .\energy
go build -trimpath -ldflags="-s -w" -o .\bin\alerts-service .\alerts
Set-Location ..
docker compose config
docker compose up -d --build
docker compose ps
Invoke-RestMethod -Uri http://localhost:8080/api/health
```

## 60. Ordre rapide pour refaire Kubernetes

Resume des commandes principales :

```powershell
Set-Location C:\docker-projets\greenops-platform
docker compose down
kubectl config use-context docker-desktop
docker build -t greenops-gateway:latest -f .\services\gateway\Dockerfile .\services
docker build -t greenops-auth-service:latest -f .\services\auth\Dockerfile .\services
docker build -t greenops-energy-service:latest -f .\services\energy\Dockerfile .\services
docker build -t greenops-alerts-service:latest -f .\services\alerts\Dockerfile .\services
kubectl apply -f .\k8s\namespace.yaml
kubectl create secret generic greenops-secrets --namespace greenops --from-literal=JWT_SECRET="GreenOpsJwtDemo2026u7N9pQ2sV5xL8rT4mK6c" --from-literal=POSTGRES_PASSWORD="GreenOpsPgDemo2026V7nR4qT9sL2" --from-literal=GRAFANA_ADMIN_PASSWORD="GreenOpsGrafanaDemo2026Q8mK5pZ1vX" --dry-run=client -o yaml | kubectl apply -f -
kubectl kustomize --load-restrictor=LoadRestrictionsNone .\k8s | kubectl apply -f -
kubectl get pods -n greenops -o wide
kubectl rollout status deployment/gateway -n greenops --timeout=180s
kubectl port-forward -n greenops svc/reverse-proxy 8080:80 --address 0.0.0.0
```

## 61. URLs finales a presenter

Sur la machine Windows :

```text
Application : http://localhost:8080
Grafana : http://localhost:3001
Prometheus : http://localhost:9090
```

Depuis une autre machine du meme reseau :

```text
Application : http://IP_DE_LA_MACHINE_WINDOWS:8080
Grafana : http://IP_DE_LA_MACHINE_WINDOWS:3001
Prometheus : http://IP_DE_LA_MACHINE_WINDOWS:9090
```

Avec hosts et Ingress si disponible :

```text
Application : http://greenops.local
Grafana : http://grafana.greenops.local
Prometheus : http://prometheus.greenops.local
```

## 62. Sources officielles consultees

- Docker Desktop Windows : https://docs.docker.com/desktop/setup/install/windows-install/
- Docker Desktop Kubernetes : https://docs.docker.com/desktop/use-desktop/kubernetes/
- Microsoft WSL : https://learn.microsoft.com/en-us/windows/wsl/install
- Git for Windows : https://git-scm.com/install/windows
- Go installation : https://go.dev/doc/install
- GitHub CLI Windows : https://github.com/cli/cli/blob/trunk/docs/install_windows.md

Ces sources ont ete verifiees le 2026-06-04 pour la partie installation sur Windows.
