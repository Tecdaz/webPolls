package main

import (
	"log"
	"net/http"
	"webpolls/db"
	"webpolls/handlers"
	"webpolls/middleware"
	"webpolls/services"
	"webpolls/utils"

	sqlc "webpolls/db/sqlc"
)

func main() {
	// inicio la conexion a la BD
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Inicializar Session Store
	utils.InitSessionStore()

	// Inyección de dependencias
	queries := sqlc.New(dbConn)

	// Inicializar servicios
	userService := services.NewUserService(queries)
	pollService := services.NewPollService(queries, dbConn)
	sseBroker := services.NewSSEBroker()

	// Inicializar handlers con los servicios
	userHandler := handlers.NewUserHandler(userService)
	pollHandler := handlers.NewPollHandler(pollService, sseBroker)
	homeHandler := handlers.NewHomeHandler(userService)

	// Crear un nuevo mux y registrar todas las rutas
	mux := http.NewServeMux()

	// Rutas de home
	mux.HandleFunc("GET /{$}", homeHandler.GetHome)

	// Rutas de usuarios
	mux.HandleFunc("POST /users/create", userHandler.CreateUser)
	// mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	// mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)
	// mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	// mux.HandleFunc("GET /users", userHandler.GetUsers) // Oculto por ahora

	// Auth routes
	mux.HandleFunc("GET /login", userHandler.GetLogin)
	mux.HandleFunc("POST /login", userHandler.PostLogin)
	mux.HandleFunc("GET /register", userHandler.GetRegister)
	mux.HandleFunc("/logout", userHandler.Logout)

	// Rutas de encuestas (Protegidas)
	mux.Handle("POST /polls/create", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.CreatePoll)))
	mux.Handle("GET /polls/{id}", middleware.OptionalAuthMiddleware(http.HandlerFunc(pollHandler.GetPollPage)))
	mux.Handle("POST /polls/{id}/vote", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.Vote)))
	mux.Handle("DELETE /polls/{id}", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.DeletePoll)))
	mux.Handle("GET /polls", middleware.OptionalAuthMiddleware(http.HandlerFunc(pollHandler.GetPolls)))
	mux.Handle("GET /my-polls", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.GetMyPolls))) // New protected route
	mux.Handle("PUT /options/{id}", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.UpdateOption)))
	mux.Handle("DELETE /polls/{poll_id}/options/{id}", middleware.AuthMiddleware(http.HandlerFunc(pollHandler.DeleteOption)))
	mux.HandleFunc("GET /polls/components/option", pollHandler.GetPollOptionInput) // Public? Used in creation form. If creation is protected, this might need to be too, but it's just a fragment.
	mux.HandleFunc("GET /events", pollHandler.SSE)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// inicio servidor
	port := ":8080"
	log.Println("Servidor corriendo en", port)
	// Usar el mux envuelto en el middleware
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
