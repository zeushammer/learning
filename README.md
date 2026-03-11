Help with packages www.alexedwards.net/blog/an-introduction-to-packages-imports-and-modules
github.com/seankhliao
github.com/TheCoreMan
cupogo.dev/episodes

github.com/eliben/code-for-blog/blog/main/2021/go-rest-servers/stdlib-middleware/stdlib-middleware.go

Players send requests
Requests enter a priority queue
Workers pull from the queue
Matches form
Backed allocats a game server using Agones

---

### 1. Install Toolchain

Install the core dependencies: **Go**, **Docker**, **Kubectl**, and **Kind**.

```bash
# Update and install build essentials
sudo apt update && sudo apt install -y build-essential curl git

# Install Go (v1.26.1)
curl -LO https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Install Docker
sudo apt install -y ca-certificates gnupg
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker $USER # Logout/Login required after this

# Install Kind & Kubectl
go install sigs.k8s.io/kind@v0.31.0
curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl && sudo mv ./kubectl /usr/local/bin/

# Install Helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh && ./get_helm.sh

```

### 2. Infrastructure Setup

Spin up the local cluster and install **Agones** to manage the game server lifecycles.

```bash
# Create the cluster
kind create cluster --name arena

# Install Agones via Helm
helm repo add agones https://agones.dev/chart/stable
helm repo update
helm install agones agones/agones --namespace agones-system --create-namespace

# Update the storage limits to run locally
kubectl set resources deployment agones-controller -n default --limits=ephemeral-storage=500Mi --requests=ephemeral-storage=100Mi
```

### 3. Build & Deploy

Build the server image and load it into the cluster nodes.

```bash
# Build the image
docker build -t arena-game-server:latest .

# Side-load into Kind
kind load docker-image arena-game-server:latest --name arena

# Deploy the manifest
kubectl apply -f deploy/agones/gameserver.yaml

```
---
