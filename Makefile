gen:
	protoc -I ./proto --go_out=./grpc_tools/pb/user_pb/  \
		--go-grpc_out=./grpc_tools/pb/user_pb/ --grpc-gateway_out=./grpc_tools/pb/user_pb/  proto/*.proto

serve:
	go run ./cmd/serve/main.go -host "127.0.0.1" -grpcPort 50051 -restPort 8080

serve2:
	go run ./cmd/serve/main.go -host "127.0.0.1" -grpcPort 50052 -restPort 8081

client:
	go run ./cmd/client/main.go

secret:
	cd grpc_tools/cert && gen.sh && cd ../..

.PHONY:
	gen run_grpc run_rest client secret run_grpc2


