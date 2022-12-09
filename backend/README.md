# Users Management Backend

## Local Development

### Using the Firebase Emulator

Build the docker image: `docker build -t epi_users_firestore -f emulator/Dockerfile.firestore .`

Run the docker container:

```bash
docker run -p 4000:4000 -p 9150:9150 -p 8090:8090 -p 9099:9099 epi_users_firestore
```

Or:

```bash 
./emulator/start_emulator.sh
```

### Run Tests

The Firebase Emulator must be running locally. Run tests with: `make test`


## BUILD
Build with docker:
```bash 
docker build -t moh_epi_auth .
```

Tag: `docker tag moh_epi_auth us-east1-docker.pkg.dev/moh-epi/app-deployments/moh_epi_auth:version`. The version format 
is: `v0.1.0`, following semver. Eg of tagging with a version:
```bash
docker tag moh_epi_auth us-east1-docker.pkg.dev/moh-epi/app-deployments/moh_epi_auth:v0.1.0
```

Push to the registry:
```bash 
docker push us-east1-docker.pkg.dev/moh-epi/app-deployments/moh_epi_auth:v0.1.0
```


## Testing
You need the firestore emulator to run our tests. Build the docker container locally:
```bash
docker build -t epi_users_firestore -f emulator/Dockerfile.firestore .
docker tag epi_users_firestore us-east1-docker.pkg.dev/moh-epi/app-deployments/epi_users_firestore:v0.1.0
```

Push the image to GCP so we can use it for testing in our automated CI:

```bash
docker push us-east1-docker.pkg.dev/moh-epi/app-deployments/epi_users_firestore:v0.1.0
```
