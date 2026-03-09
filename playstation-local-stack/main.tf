# 1. Create the Local Kubernetes Cluster (Kind)
resource "kind_cluster" "default" {
  name = "playstation-local"
  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"
    node {
      role = "control-plane"
      extra_port_mappings {
        container_port = 30000
        host_port      = 7654
        protocol       = "UDP"
      }
    }
  }
}

# 2. Install Agones via Helm
resource "helm_release" "agones" {
  name             = "agones"
  repository       = "https://agones.dev/chart/stable"
  chart            = "agones"
  namespace        = "agones-system"
  create_namespace = true
  depends_on       = [kind_cluster.default]

  set {
    name  = "agones.allocator.http.enabled"
    value = "false"
  }
}

# 3. Install Crossplane via Helm
resource "helm_release" "crossplane" {
  name             = "crossplane"
  repository       = "https://charts.crossplane.io/stable"
  chart            = "crossplane"
  namespace        = "crossplane-system"
  create_namespace = true
  depends_on       = [kind_cluster.default]
}

# 4. Define the Agones Fleet (The "Game Server")
resource "kubernetes_manifest" "game_fleet" {
  manifest = {
    apiVersion = "agones.dev/v1"
    kind       = "Fleet"
    metadata = {
      name      = "terraform-fleet"
      namespace = "default"
    }
    spec = {
      replicas = 2
      template = {
        spec = {
          ports = [{ name = "default", containerPort = 7654 }]
          template = {
            spec = {
              containers = [{
                name  = "game-server"
                image = "us-docker.pkg.dev/agones-images/udp-server:0.32"
              }]
            }
          }
        }
      }
    }
  }
  depends_on = [helm_release.agones]
}
