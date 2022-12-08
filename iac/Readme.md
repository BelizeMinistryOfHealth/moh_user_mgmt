# IAC 

## Authentication
Before running any terraform command, you need to authenticate with GCP:

```bash
gcloud auth application-default login --no-launch-browser
```

## Starting from Scratch
- Initialize terraform: `terraform init --reconfigure`. This will store the terraform state in GCS.
