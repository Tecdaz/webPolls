# WebPolls – Instrucciones de ejecución

Aplicación backend en Go con base de datos PostgreSQL. Incluye un esquema SQL que se aplica al iniciar para crear tablas si no existen.

## Requisitos
- Docker y Docker Compose v2
- Go (solo si deseas ejecutarlo localmente sin Docker)

## Estructura relevante
- `main.go`: servidor HTTP y conexión a PostgreSQL; aplica el esquema al iniciar.
- `db/schema/schema.sql`: definición de tablas e índices.
- `Dockerfile`: build multi-stage para generar una imagen mínima del backend.
- `docker-compose.yml`: orquesta `postgres` y `backend`.
- `.env`: credenciales y configuración (usado por Compose y opcionalmente por `godotenv` en local).
- `static/`: archivos estáticos servidos en `GET /`.

## Variables de entorno
Archivo `.env` esperado:
```
DB_USER=usuario (cambiar)
DB_PASSWORD=contraseña (cambiar)
DB_NAME=nombre_de_la_base_de_datos (cambiar)
```

Para ejecución local (sin Docker), usa:
```
DB_HOST=localhost
```

---

## Ejecución con Docker Compose (recomendado)

1) Construir y levantar servicios
```bash
docker compose build
docker compose up -d
```

2) Probar el backend
```bash
curl http://localhost:8080/
```

3) Apagar
```bash
docker compose down
```

Notas:
- El contenedor del backend copia `db/schema/schema.sql` dentro de la imagen y lo ejecuta al arrancar para crear tablas.
- El servicio `postgres` usa las credenciales del `.env`. En Compose el backend se conecta con `DB_HOST=postgres` y `DB_PORT=5432`.

---

## Ejecución local (sin Docker)

1) Asegúrate de tener PostgreSQL accesible en `localhost:5432`.

2) Ejecuta la app
```bash
go mod init webpolls.com/webpolls
go mod tidy
go run main.go
```

La app:
- Conecta a PostgreSQL con las variables anteriores.
- Verifica la conexión (`Ping`).
- Aplica el esquema `db/schema/schema.sql` (crea tablas/índices si no existen).
- Arranca en `:8080`.

---

## Endpoints
- `GET /` → sirve `static/index.html`.

---

## Comandos útiles
```bash
# Reconstruir solo el backend
docker compose build --no-cache backend

# Logs
docker compose logs -f backend

# Probar el puerto localmente
curl -v http://localhost:8080/
```