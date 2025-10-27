# -----------------
# Pruebas API de Encuestas (Polls)
# -----------------

# 1. Crear una nueva encuesta y capturar su ID
POST http://localhost:8080/users/create
Content-Type: application/json
```json
{
  "username": "testuser",
  "email": "testuser@example.com",
  "password": "testpassword"
}
```
HTTP 201
[Captures]
new_user_id: jsonpath "$.data.id"

POST http://localhost:8080/polls/create
Content-Type: application/json
```json
{
  "question": "¿Cuál es tu lenguaje favorito?",
  "options": [
    { "content": "Go" },
    { "content": "Python" },
    { "content": "Rust" }
  ],
  "user_id": {{new_user_id}}
}
```
HTTP 201
[Captures]
new_poll_id: jsonpath "$.data.id"
[Asserts]
jsonpath "$.message" == "Encuesta creada correctamente"
jsonpath "$.data.title" == "¿Cuál es tu lenguaje favorito?"
jsonpath "$.data.options[0].content" == "Go"
jsonpath "$.data.options[1].content" == "Python"
jsonpath "$.data.options[2].content" == "Rust"
jsonpath "$.data.user_id" == {{new_user_id}}

# 2. Obtener la encuesta recién creada
# Nota: El endpoint de encuestas usa un query param `id`
GET http://localhost:8080/polls/{{new_poll_id}}
Accept: application/json

HTTP 200
[Asserts]
jsonpath "$.message" == "Encuesta obtenida correctamente"
jsonpath "$.data.id" == {{new_poll_id}}
jsonpath "$.data.title" == "¿Cuál es tu lenguaje favorito?"
jsonpath "$.data.options[0].content" == "Go"
jsonpath "$.data.options[1].content" == "Python"
jsonpath "$.data.options[2].content" == "Rust"
jsonpath "$.data.user_id" == {{new_user_id}}

# 3. Eliminar la encuesta
DELETE http://localhost:8080/polls/{{new_poll_id}}
Accept: application/json

HTTP 200
[Asserts]
jsonpath "$.message" == "Encuesta eliminada correctamente"
jsonpath "$.data" == null

# 4. Eliminar el usuario
DELETE http://localhost:8080/users/{{new_user_id}}
Accept: application/json

HTTP 200
[Asserts]
jsonpath "$.message" == "Usuario eliminado correctamente"

POST http://localhost:8080/users/create
Content-Type: application/json
```json
{
  "username": "agus",
  "email": "agus2@gmail.com",
  "password": "123456"
}
```
