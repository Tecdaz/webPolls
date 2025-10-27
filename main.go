package main

import (
	"log"
	"net/http"
	"webpolls/db"
	"webpolls/handlers"
	"webpolls/middleware"
	"webpolls/services"

	sqlc "webpolls/db/sqlc"
)

func main() {
	// inicio la conexion a la BD
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Inyección de dependencias
	queries := sqlc.New(dbConn)

	// Inicializar servicios
	userService := services.NewUserService(queries)
	pollService := services.NewPollService(queries, dbConn)

	// Inicializar handlers con los servicios
	userHandler := handlers.NewUserHandler(userService)
	pollHandler := handlers.NewPollHandler(pollService)

	// Crear un nuevo mux y registrar todas las rutas
	mux := http.NewServeMux()

	// Rutas de usuarios
	mux.HandleFunc("POST /users/create", userHandler.CreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("GET /users", userHandler.GetUsers)

	// Rutas de encuestas
	mux.HandleFunc("POST /polls/create", pollHandler.CreatePoll)
	mux.HandleFunc("GET /polls/{id}", pollHandler.GetPoll)
	mux.HandleFunc("DELETE /polls/{id}", pollHandler.DeletePoll)
	mux.HandleFunc("GET /polls", pollHandler.GetPolls)
	mux.HandleFunc("PUT /options/{id}", pollHandler.UpdateOption)
	mux.HandleFunc("DELETE /polls/{poll_id}/options/{id}", pollHandler.DeleteOption)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Aplicar el middleware a todo el mux
	loggedMux := middleware.LoggingMiddleware(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// inicio servidor
	port := ":8080"
	log.Println("Servidor corriendo en", port)
	// Usar el mux envuelto en el middleware
	if err := http.ListenAndServe(port, loggedMux); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
