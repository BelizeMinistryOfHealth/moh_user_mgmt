FROM golang:1.19-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/server cmd/server/main.go
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /bin/server /bin/server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/bin/server"]
