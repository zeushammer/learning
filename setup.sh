#!/bin/bash
set -e

echo "🚀 Starting Environment Setup for Arena..."

# 1. System Dependencies
sudo apt update && sudo apt install -y build-essential curl git ca-certificates gnupg

# 2. Go Installation (v1.26.1)
if ! command -v go &> /dev/null; then
    echo "Installing Go v1.26.1..."
    curl -LO https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
    rm go1.26.1.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
fi

# 3. Docker Installation
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io
    sudo usermod -aG docker $USER
    echo "⚠️ NOTE: You may need to log out and back in for Docker group changes to take effect."
fi

# 4. Kubernetes Tools (Kind, Kubectl, Helm)
echo "Installing K8s toolchain..."
go install sigs.k8s.io/kind@v0.31.0

if ! command -v kubectl &> /dev/null; then
    K8S_VER=$(curl -L -s https://dl.k8s.io/release/stable.txt)
    curl -LO "https://dl.k8s.io/release/${K8S_VER}/bin/linux/amd64/kubectl"
    chmod +x ./kubectl && sudo mv ./kubectl /usr/local/bin/
fi

if ! command -v helm &> /dev/null; then
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 get_helm.sh && ./get_helm.sh && rm get_helm.sh
fi

# 5. Cluster & Agones Initialization

echo "Initializing Kind Cluster 'arena'..."
kind create cluster --name arena || echo "Cluster already exists."

echo "Installing Agones..."
helm repo add agones https://agones.dev/chart/stable
helm repo update
helm install agones agones/agones --namespace agones-system --create-namespace --wait

echo "✅ Setup Complete! Run 'kubectl get pods -A' to verify."
