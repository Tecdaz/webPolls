package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"webpolls/handlers"

	sqlc "webpolls/db/sqlc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// cargar variables .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema.")
	}

	// conexion a BD
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable",
		dbUser, dbPassword, dbName, dbHost,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil { //ping para verificar conexion
		log.Fatal("Error pinging database:", err)
	}

	// aplicar schema.sql
	schemaBytes, err := os.ReadFile("db/schema/schema.sql")
	if err != nil {
		log.Fatalf("Error leyendo schema.sql: %v", err)
	}
	if _, err := db.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Error aplicando schema: %v", err)
	}

	fmt.Println("Conexión a la base de datos exitosa")

	// inicio handlers
	queries := sqlc.New(db)
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
