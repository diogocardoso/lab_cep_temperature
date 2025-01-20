# Etapa 1: Construir a aplicação
FROM golang:1.21.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY config.yaml ./config.yaml 

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun cmd/main.go

# Etapa 2: Criar a imagem final
FROM alpine:latest
COPY --from=builder /app/cloudrun .
COPY --from=builder /app/config.yaml ./config.yaml
CMD ["./cloudrun"]