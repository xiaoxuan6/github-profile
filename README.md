# github-profile

## Start locally

- go 1.22 with go mod
- go run main.go
- with token GITHUB_TOKE=xxxx go run main.go

## Docker

#### environment

```docker
docker run --name=github-profile -e GITHUB_TOKEN="xxx" -p 8080:8080 ghcr.io/xiaoxuan6/github-profile:latest 
```

#### .env file

```docker
docker run --name=github-profile -v $(pwd)/.env:/src/.env -p 8080:8080 ghcr.io/xiaoxuan6/github-profile:latest 
```