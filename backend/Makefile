API_KEY?=""
GOOGLE_APPLICATION_CREDENTIALS?=""
FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"

tidy:
	go mod tidy

clean:
	rm -rf ./bin Gopkg.lock

test:
	PROJECT_ID=demo-project FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 API_KEY=API_KEY FIRESTORE_EMULATOR_HOST=localhost:8090 GOOGLE_APPLICATION_CREDENTIALS=GOOGLE_APPLICATION_CREDENTIALS go test -coverprofile=coverage.out ./...

# Build
build-linux:
	export GO111MODULE=on
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	chmod +x bin/server

build-macos:
	export GO111MODULE=on
	env GOOS=darwin go build -o bin/server cmd/server/main.go
	chmod +x bin/server

build-user-cli:
	export GO111MODULE=on
	env GOOS=darwin go build -o bin/cli cmd/cli/main.go
	chmod +x bin/cli
