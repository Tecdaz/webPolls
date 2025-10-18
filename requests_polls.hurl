# pruebas API de encuestas
#create
POST http://localhost:8080/polls/create
Content-Type: application/json

{
  "question": "¿Cuál es tu lenguaje favorito?",
  "options": ["Go", "Python", "Rust"]
}

HTTP/1.1 201
[Asserts]
  jsonpath "$.poll.title" == "¿Cuál es tu lenguaje favorito?"
  jsonpath "$.options[0].content" == "Go"
  jsonpath "$.options[1].content" == "Python"


