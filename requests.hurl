# Pruebas API usuarios

# Crear usuario
POST http://localhost:8080/users/create
Content-Type: application/json

{
  "username": "agustina",
  "email": "agus@gmail.com",
  "password": "123456"
}

HTTP/1.1 200
[Asserts]
  jsonpath "$.user.username" == "agustina"


# Obtener usuario por id
GET http://localhost:8080/users/5
Accept: application/json

HTTP/1.1 200
[Asserts]
  jsonpath "$.user.username" == "agustina"

# Actualizar solo email de usuario
PUT http://localhost:8080/users/5
Content-Type: application/json

{
  "username": "agustina",
  "email": "modificado@gmail.com"
}

HTTP/1.1 200
[Asserts]
  jsonpath "$.user.email" == "modificado@gmail.com"

# Actualizar solo username de usuario
PUT http://localhost:8080/users/5
Content-Type: application/json

{
  "username": "modificado",
  "email": "modificado@gmail.com"
}

HTTP/1.1 200
[Asserts]
  jsonpath "$.user.username" == "modificado"

# Eliminar usuario
DELETE http://localhost:8080/users/5
Accept: application/json

HTTP/1.1 200
[Asserts]
  jsonpath "$.message" == "User deleted successfully"
