gen:
	protoc -I ./proto --go_out=../utools/pb/userPB/  --go-grpc_out=../utools/pb/userPB/ --grpc-gateway_out=../utools/pb/userPB/  proto/*.proto

run_grpc:
	go run ./cmd/serve/main.go -port 50051

run_rest:
	go run ./cmd/serve/main.go -type rest -port 8080 -endpoint 0.0.0.0:50051

client:
	go run ./cmd/client/main.go

secret:
	cd cert && gen.sh && cd ..

.PHONY:
	gen run_grpc run_rest client secret


