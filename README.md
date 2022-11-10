# grpc_user_service

## 1. 安装protoc,下载对应系统的版本

```text
https://github.com/protocolbuffers/protobuf/releases
```

## 2. 安装protoc-gen-go和protoc-gen-go-grpc

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest //安装grpc-gateway
```

## 3. 下载子模块utools工具包

```text
git submodule init
git submodule update
```

## 4. 项目启动

- 本地创建pgsql, 创建user数据库

```shell
docker run --name mypostgres -d -p 5432:5432 -e POSTGRES_PASSWORD=123456 postgres
```

- 启动etcd注册中心

```shell
docker run -d -p 12379:2379 -p 12380:2380 -v etcd_data:/etcd-data/member \
 --name exam-etcd quay.io/coreos/etcd:latest  usr/local/bin/etcd \
  --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379 \
   --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 \
     --initial-advertise-peer-urls http://0.0.0.0:2380  --initial-cluster s1=http://0.0.0.0:2380 \
      --initial-cluster--token tkn  --initial-cluster-state new
```

- make serve 启动grpc和rest服务（make 需提前安装 无make使用go run ./cmd/serve/main.go -host "127.0.0.1" -grpcPort 50051 -restPort 8080 启动）

