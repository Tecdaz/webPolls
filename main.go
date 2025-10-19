package main

import (
	"log"
	"net/http"
	"webpolls/db"
	"webpolls/handlers"

	sqlc "webpolls/db/sqlc"
)

func main() {
	// inicio la conexion a la BD
	dbConn := db.InitDB()
	defer dbConn.Close()

	// inicio handlers
	queries := sqlc.New(dbConn)
	userHandler := handlers.NewUserHandler(queries)
	pollHandler := handlers.NewPollHandler(queries)

	// Rutas de usuarios
	http.HandleFunc("POST /users/create", userHandler.CreateUser)
	http.HandleFunc("GET /users/{id}", userHandler.GetUser)
	http.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)
	http.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)

	// Rutas de encuestas
	http.HandleFunc("POST /polls/create", pollHandler.CreatePoll)
	http.HandleFunc("GET /polls/", pollHandler.GetPoll)
	http.HandleFunc("DELETE /polls/", pollHandler.DeletePoll)

	// ruta principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))) //para servir el contenido de static/

	// inicio servidor
	port := ":8080"
	log.Println("Servidor corriendo en", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
