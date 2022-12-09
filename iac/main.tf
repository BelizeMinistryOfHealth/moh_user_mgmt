provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

locals {
  users_mgmt_container = "us-east1-docker.pkg.dev/moh-epi/app-deployments/moh_epi_auth:v0.2.0"
}

resource "google_artifact_registry_repository" "my-repo" {
  location      = var.region
  repository_id = "app-deployments"
  description   = "docker repository for application deployments"
  format        = "DOCKER"
}

resource "google_cloud_run_service" "users-mgmt" {
  name     = "users-mgmt"
  location = var.region

  template {
    spec {
      containers {
        image = local.users_mgmt_container
        env {
          name = "PROJECT_ID"
          value = var.project_id
        }
        env {
          name = "API_KEY"
          value = var.firebase_api_key
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "allUsers" {
  service  = google_cloud_run_service.users-mgmt.name
  location = google_cloud_run_service.users-mgmt.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

terraform {
  backend "gcs" {
    bucket = "moh-epi-tfstates"
  }
}
