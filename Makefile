gen:
	protoc -I ./proto --go_out=../tools/pb/user_pb/  --go-grpc_out=../tools/pb/user_pb/ --grpc-gateway_out=../tools/pb/user_pb/  proto/*.proto

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


