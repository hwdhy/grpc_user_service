gen:
	protoc --go_out=./pb  --go-grpc_out=./pb  proto/*.proto

run:
	go run ./cmd/serve/main.go

client:
	go run ./cmd/client/main.go

.PHONY:
	gen run client