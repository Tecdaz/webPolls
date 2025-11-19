# compilacion
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copiar el código fuente
COPY . .

RUN go mod download

# Compilar binario
RUN go build -o ./bin/webpolls ./main.go

## ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app

# Copiar binario y recursos necesarios
COPY --from=builder /app/bin/webpolls /app/webpolls
COPY --from=builder /app/static/ /app/static/
COPY --from=builder /app/db/schema/schema.sql /app/db/schema/schema.sql

# Instalar cliente de postgres para el entrypoint
RUN apk add --no-cache postgresql-client

# Copiar y configurar entrypoint
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["./webpolls"]