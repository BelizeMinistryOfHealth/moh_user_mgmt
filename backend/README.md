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


