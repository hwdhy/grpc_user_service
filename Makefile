gen:
	protoc --go_out=./pb  --go-grpc_out=./pb  proto/*.proto

serve:
	go run ./cmd/serve/main.go

.PHONY:
	gen serve