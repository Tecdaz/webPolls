# WebPolls

Aplicación web simple para crear y visualizar encuestas con opciones de respuesta de tipo multiple choice.

La app actualmente expone una página estática en `http://localhost:8080/` servida por el servidor en Go (`main.go`). El objetivo del dominio es permitir que un usuario cree encuestas (polls) y que cada encuesta tenga varias opciones (options) entre las cuales se pueda elegir.

## Modelos de datos

A continuación se describe la información almacenada para cada modelo según el diagrama provisto (`user` — `poll` — `option`).

### user
- **id**: `serial` (PK)
  - Identificador único del usuario.
- **username**: `varchar(20)`
  - Nombre de usuario visible y único (recomendado aplicar restricción de unicidad).
- **password**: `varchar`
  - Hash de la contraseña del usuario. No se deben almacenar contraseñas en texto plano.
- **email**: `varchar(20)`
  - Email del usuario.

### poll
- **id**: `serial` (PK)
  - Identificador único de la encuesta.
- **title**: `varchar(200)`
  - Título o pregunta principal de la encuesta.
- **user_id**: `int` (FK → `user.id`)
  - Usuario propietario/creador de la encuesta.

Relación: Un `user` puede tener muchas `poll` (1:N).

### option
- **id**: `serial` (PK)
  - Identificador único de la opción.
- **content**: `varchar(50)`
  - Texto de la opción que verá el usuario al votar.
- **poll_id**: `int` (PK, FK → `poll.id`)
  - Encuesta a la que pertenece la opción. En el diagrama figura como parte de la clave (PK, FK), lo que sugiere una clave compuesta (`id`, `poll_id`). Alternativamente, puede modelarse como PK simple en `id` y `poll_id` como FK con índice.
- **correct**: `boolean`
  - Marca si la opción es la correcta (útil si la encuesta funciona como cuestionario). Para encuestas sin respuesta correcta, puede ignorarse o dejarse en `false`.

Relación: Una `poll` tiene muchas `option` (1:N).

## Estructura del proyecto

```
webPolls/
├─ main.go              # Servidor HTTP (Go)
├─ go.mod               # Módulo de Go
├─ DockerFile         
├─ Makefile           
└─ static/              # Archivos estáticos (frontend)
   ├─ styles.css
   └─ logo.svg (opcional)
├─ db/                  # Base de datos
   ├─ schema/           # Esquema de la base de datos
   │  └─ schema.sql     # Esquema de la base de datos
   └─ queries/          # Consultas a la base de datos
      ├─ users.sql     # Consultas a la tabla de usuarios
      ├─ polls.sql     # Consultas a la tabla de encuestas
      └─ options.sql   # Consultas a la tabla de opciones
├─ sqlc/               # Generación de código (sqlc)
   ├─ users.sql.go    # Código generado para la tabla de usuarios
   ├─ polls.sql.go    # Código generado para la tabla de encuestas
   └─ options.sql.go  # Código generado para la tabla de opciones
├─ utils/              # Utilidades
├─ components/         # Componentes reutilizables        
├─ views/              # Paginas con templ
   ├─ home.templ
   ├─ index.templ 
   ├─ polls.templ 
   ├─ users.templ   
├─ handlers/            # Handlers de rutas
   ├─ homeHandler.go
   ├─ pollsHandler.go 
   ├─ userHandler.go
├─ services/            # Logica de negocio
   ├─ userHandler.go
   ├─ userHandler.go 
```

## Ejecución local

### Requisitos
- Docker y Docker Compose
- Make
- Go 1.20+ (opcional, para desarrollo local)
- sqlc (opcional, para regenerar código)

### Comandos disponibles

Para ver todos los comandos disponibles:
```bash
make help
```

### Ejecución completa con Docker

1. **Ejecutar la aplicación completa** (recomendado):
   ```bash
   make run
   ```
   Este comando limpia contenedores anteriores, construye la imagen, levanta los servicios y crea datos de prueba.

2. **Abrir en el navegador**:
   - `http://localhost:8080/`

### Comandos útiles

- **Instalar dependencias de Go**:
  ```bash
  make install
  ```

- **Construir solo la imagen Docker**:
  ```bash
  make docker-build
  ```

- **Levantar contenedores**:
  ```bash
  make docker-up
  ```

- **Bajar contenedores**:
  ```bash
  make docker-down
  ```

- **Ver logs de la aplicación**:
  ```bash
  make logs
  ```

- **Reiniciar la aplicación**:
  ```bash
  make restart
  ```

- **Generar código SQLC**:
  ```bash
  make sqlc
  ```

- **Crear datos de prueba**:
  ```bash
  make seed
  ```

### Ejecución local sin Docker

Para desarrollo local (requiere base de datos PostgreSQL ejecutándose):
```bash
make run-local
```


## Integrantes del grupo

- Agustina Pereyra
- Joaquin Loza Ciappa
- Santiago Arias Ocampo