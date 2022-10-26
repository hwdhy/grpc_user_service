FROM golang:latest as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app ./cmd/serve/main.go


FROM alpine as prod

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=build /go/release /app/

CMD ["/app/app -port 50051"]
CMD ["/app/app -type rest -port 8080 -endpoint 0.0.0.0:50051"]