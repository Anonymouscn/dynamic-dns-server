FROM golang:1.21-alpine AS builder
LABEL maintainer="Anonymous"

ARG LDFLAGS_NAME="-X main.name="
ARG LDFLAGS_VERSION="-X main.version="
ARG ENTRY="./cmd"
ARG PRODUCT="./bin"

ENV NAME="dynamic-dns-server"
ENV VERSION="1.0"

WORKDIR /app

COPY . .
RUN go mod download \
    && go clean \
    && go build -ldflags "${LDFLAGS_NAME}${NAME} ${LDFLAGS_VERSION}${VERSION}" -o ${PRODUCT}/start ${ENTRY}

FROM alpine:3.7
LABEL maintainer="Anonymous"

COPY --from=builder /app/bin/start /app/
COPY application.yml /app/

ENV TZ=Asia/Shanghai
RUN echo -e https://mirrors.aliyun.com/alpine/v3.7/main/ > /etc/apk/repositories && \
    apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && chmod +x /app/start

WORKDIR /app

ENTRYPOINT ["./start"]