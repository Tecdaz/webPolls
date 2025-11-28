# compilacion
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copiar el código fuente
# Instalar Node.js y pnpm para construir CSS
RUN apk add --no-cache nodejs npm
RUN npm install -g pnpm

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Copiar el código fuente
COPY . .

# Construir CSS
RUN pnpm build:css

# Compilar binario
RUN go build -o ./bin/webpolls ./main.go

## ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app

# Copiar binario y recursos necesarios
COPY --from=builder /app/bin/webpolls /app/webpolls
COPY --from=builder /app/static/ /app/static/
COPY --from=builder /app/db/schema/schema.sql /app/db/schema/schema.sql

EXPOSE 8080

CMD ["./webpolls"]