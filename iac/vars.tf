variable "project_id" {
  type        = string
  description = "The Project ID"
}

variable "region" {
  type        = string
  description = "GCP Region"
}


variable "zone" {
  type        = string
  description = "GCP Zone"
}

variable "tfstate" {
  type        = string
  description = "Holds the terraform state"

}

variable "firebase_api_key" {
  type        = string
  description = "API Key for interacting with Firebase"
}
