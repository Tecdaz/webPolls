# WebPolls

- [Cómo empezar](#-cómo-empezar)
- [Estructura del proyecto](#-estructura-del-proyecto)
- [Desarrollo y ejecucion](#-desarrollo-y-ejecucion)
- [Modelos de datos](#-modelos-de-datos)
- [Autores](#-autores)

Aplicación web simple para crear y visualizar encuestas con opciones de votación de tipo multiple choice.

## Cómo empezar

### Script de inicio

Levanta docker, construye la imagen y popula la base de datos.

Requisitos: Docker, Docker Compose y Make

1. Ejecutar el comando principal:
   ```bash
   make run
   ```

El servidor quedará corriendo en el puerto 8080 y la base de datos en el puerto 5432.


## Estructura del proyecto

```
├───docker-compose.yml
├───Dockerfile
├───Makefile            # Archivo de comandos para facilitar tareas
├───go.mod
├───go.sum
├───.env                # Variables de entorno
├───main.go
├───docker-entrypoint.sh             # Script para popular la base de datos
├───sqlc.yaml           # Configuración de sqlc
├───components/         # Componentes reutilizables (Templ)
├───db/
│   ├───connection.go   # Conexion a la base de datos
│   ├───queries/        # Queries utilizadas en la capa de servicio
│   ├───schema/         # Esquema de la base de datos
│   └───sqlc/           # Archivos generados por sqlc
├───handlers/           # Capa de presentación de la api
├───middleware/         # Middleware para logging
├───services/           # Capa de negocio
├───static/             # Archivos estaticos
├───tests/              # Archivos de tests
├───utils/              # Funciones utilitarias
└───views/              # Plantillas de html
```

## Modelos de datos

A continuación se describe la información almacenada para cada modelo según el diagrama provisto (`user` — `poll` — `option` - `result`).

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

Relación: Una `poll` tiene muchas `option` (1:N).

### result
- **id**: `serial` (PK)
  - Identificador único del resultado.
- **option_id**: `int` (FK → `option.id`)
  - Opción seleccionada por el usuario.
- **poll_id**: `int` (FK → `poll.id`)
  - Encuesta a la que pertenece el resultado.
- **user_id**: `int` (FK → `user.id`)
  - Usuario que realizó la votacion.

A parte de un id unico que identifica cada resultado, hay una clave compuesta (option_id, poll_id, user_id) que identifica cada resultado.

Relación: Un `result` pertenece a una `option`, una `poll` y un `user` (1:N).

## Desarrollo y ejecucion

Agrupamos los comandos utiles para la etapa de desarrollo y ejecucion en el archivo `Makefile`
con el objetivo de automatizar tareas repetitivas y ofrecer un mejor experiencia al usuario
y sistemas que automatizan operaciones.


## Frontend

El frontend de la aplicación está construido utilizando **Templ**, una librería de Go para generar HTML de manera eficiente y tipada.

La aplicacion es reactiva usando unicamente htmx.

Para acceder al frontend se debe acceder mediante el navegador a la direccion `http://localhost:8080`.

### Estructura del Frontend

- **components/**: Contiene componentes reutilizables de la interfaz (ej: barra de navegación).
- **views/**: Contiene las páginas principales de la aplicación.

### Secciones

Actualmente tenemos tres secciones principales:
- **Inicio**: Página de bienvenida en la ruta `/`.
- **Encuestas**: Muestra todas las encuestas disponibles en la ruta `/polls`. Permite crear, eliminar y ver encuestas.
- **Usuarios**: Muestra todos los usuarios disponibles en la ruta `/users`. Permite crear y eliminar usuarios.

Nota: Actualmente cada poll se crea con un usuario hardcodeado para fines de demostración, pendiente de implementar autenticación completa.


##  Autores

- Agustina Pereyra
- Joaquin Loza Ciappa
- Santiago Arias Ocampo