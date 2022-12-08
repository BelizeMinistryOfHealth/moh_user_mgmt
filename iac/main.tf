provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

resource "google_artifact_registry_repository" "my-repo" {
  location      = var.region
  repository_id = "app-deployments"
  description   = "docker repository for application deployments"
  format        = "DOCKER"
}

terraform {
  backend "gcs" {
    bucket = "moh-epi-tfstates"
  }
}
