# compilacion
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copiar módulos primero para cachear dependencias si existen
COPY go.mod go.sum ./
RUN [ -f go.mod ] || go mod init webpolls.com/webpolls
RUN go mod tidy || true
RUN go mod download || true

# Copiar el código fuente
COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN sqlc generate

# Compilar binario
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./

## ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app

# Copiar binario y recursos necesarios
COPY --from=builder /app/main /app/main
COPY --from=builder /app/static/ /app/static/
COPY --from=builder /app/db/schema/schema.sql /app/db/schema/schema.sql

EXPOSE 8080

ENTRYPOINT ["./main"]