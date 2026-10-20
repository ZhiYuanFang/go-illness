# 多阶段构建：CI / 本地均可从源码产出 illness-service 镜像（不依赖预编译 main）。
FROM m.daocloud.io/docker.io/library/golang:1.19-alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/illness-service .

FROM m.daocloud.io/docker.io/library/alpine:3.19

RUN apk add --no-cache tzdata ca-certificates wget
ENV TZ=Asia/Shanghai

ENV WORKDIR=/app
ENV GF_GCFG_FILE=manifest/config/config.yaml
WORKDIR ${WORKDIR}

COPY --from=builder /out/illness-service ${WORKDIR}/illness-service
COPY manifest/config/config.yaml ${WORKDIR}/manifest/config/config.yaml
COPY resource ${WORKDIR}/resource

EXPOSE 9686
HEALTHCHECK --interval=10s --timeout=3s --retries=6 CMD wget -q -O - http://127.0.0.1:9686/api.json > /dev/null || exit 1
CMD ["./illness-service"]
