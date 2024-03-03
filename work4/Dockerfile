FROM golang:1.20 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .
RUN --mount=type=cache,mode=0777,id=go-mode,target=/go/pkg/mod \
    go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags="-w -s" -o ./main
WORKDIR /app
RUN mkdir publish  \
    && cp main publish  \
    && cp -r config publish

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/publish .

# 指定运行时环境变量
EXPOSE 8000

ENTRYPOINT ["./main"]