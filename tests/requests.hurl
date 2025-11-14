# -----------------
# Pruebas API de Usuarios (con datos originales)
# -----------------

# 1. Crear un nuevo usuario y capturar su ID
POST http://localhost:8080/users/create
Content-Type: application/json
```json
{
  "username": "agustina",
  "email": "agus@gmail.com",
  "password": "123456"
}
```
HTTP 201
[Captures]
user_id: jsonpath "$.data.id"
[Asserts]
jsonpath "$.message" == "Usuario creado correctamente"
jsonpath "$.data.username" == "agustina"

# 2. Obtener el usuario recién creado
GET http://localhost:8080/users/{{user_id}}
Accept: application/json

HTTP 200
[Asserts]
jsonpath "$.message" == "Usuario obtenido correctamente"
jsonpath "$.data.id" == {{user_id}}
jsonpath "$.data.username" == "agustina"

# 3. Actualizar solo el email del usuario
PUT http://localhost:8080/users/{{user_id}}
Content-Type: application/json
```json
{
  "email": "modificado@gmail.com"
}
```
HTTP 200
[Asserts]
jsonpath "$.message" == "Usuario actualizado correctamente"
jsonpath "$.data.email" == "modificado@gmail.com"
jsonpath "$.data.username" == "agustina"

# 4. Actualizar solo el nombre del usuario
PUT http://localhost:8080/users/{{user_id}}
Content-Type: application/json
```json
{
  "username": "agustina_modificada"
}
```
HTTP 200
[Asserts]
jsonpath "$.message" == "Usuario actualizado correctamente"
jsonpath "$.data.username" == "agustina_modificada"
jsonpath "$.data.email" == "modificado@gmail.com"

# 5. Eliminar el usuario
DELETE http://localhost:8080/users/{{user_id}}
Accept: application/json

HTTP 200
[Asserts]
jsonpath "$.message" == "Usuario eliminado correctamente"
jsonpath "$.data.username" == "agustina_modificada"