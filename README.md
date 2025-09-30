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
└─ static/              # Archivos estáticos (frontend)
   ├─ index.html
   ├─ styles.css
   └─ logo.svg (opcional)
```

## Ejecución local

Requisitos: Go 1.20+ (o compatible)

1. Iniciar el modulo:
   ```bash
   go mod init webPolls/webPolls.com
   ```
2. Instalar dependencias del módulo:
   ```bash
   go mod tidy
   ```
3. Ejecutar el servidor:
   ```bash
   go run main.go
   ```
4. Abrir en el navegador:
   - `http://localhost:8080/`


## Integrantes del grupo

- Agustina Pereyra
- Joaquin Loza Ciappa
- Santiago Arias Ocampo