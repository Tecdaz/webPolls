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

	// defino las rutas de users
	http.HandleFunc("/users/create", userHandler.CreateUser)
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// defino las rutas de polls
	http.HandleFunc("/polls/create", pollHandler.CreatePoll)
	http.HandleFunc("/polls/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pollHandler.GetPoll(w, r)
		case http.MethodDelete:
			pollHandler.DeletePoll(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

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
