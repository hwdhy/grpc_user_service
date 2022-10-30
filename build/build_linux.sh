rm app

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o app ../cmd/serve/main.go