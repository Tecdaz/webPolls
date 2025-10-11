# compilacion
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copiar el código fuente
COPY . .

RUN go mod download

# Compilar binario
RUN sh build.sh

## ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app

# Copiar binario y recursos necesarios
COPY --from=builder /app/bin/webpolls /app/webpolls
COPY --from=builder /app/static/ /app/static/
COPY --from=builder /app/db/schema/schema.sql /app/db/schema/schema.sql

EXPOSE 8080

ENTRYPOINT ["./webpolls"]