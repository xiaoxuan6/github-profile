FROM golang:1.22.5-alpine3.20 as builder

WORKDIR /go/app/src

COPY . .
RUN apk add --no-cache upx || \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -ldflags="-s -w" -o github_profile main.go && \
    [ -e /usr/bin/upx ] && upx github_profile || echo

FROM alpine:latest

WORKDIR /src

COPY .env.example entrypoint.sh ./
COPY static ./static
COPY templates ./templates
COPY --from=builder /go/app/src/github_profile .

RUN cp .env.example .env && \
    chmod +x entrypoint.sh && \
    chmod +x .env

ENV GITHUB_TOKEN=""

EXPOSE 8080

ENTRYPOINT ["sh", "./entrypoint.sh"]
