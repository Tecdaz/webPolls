package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"github.com/joho/godotenv"
	sqlc "webpolls.com/webpolls/db/sqlc"
	_ "github.com/lib/pq"
)

// Crea tablas a partir del esquema en disco en el arranque

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveIndex)

	port := ":8080"
	_ = godotenv.Load()

	//conexion a la base de datos
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable",
		dbUser, dbPassword, dbName, dbHost)
	db, errDb := sql.Open("postgres", connStr)

	if errDb != nil {
		fmt.Println("Error connecting to the database:", errDb)
		return
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	// Verificar la conexión
	errDb = db.Ping()
	if errDb != nil {
		log.Fatal("Error pinging database:", errDb)
	}

	// Ejecutar el esquema para asegurar que existan las tablas (idempotente por IF NOT EXISTS)
	schemaBytes, err := os.ReadFile("db/schema/schema.sql")
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}
	if _, err := db.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Error applying schema: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	// Crear usuario si no existe (buscar por username primero)
	username := "testuser"
	password := "password123"
	email := "jklsajdlaskj@skd.com"

	userRow, err := queries.GetUserByUsername(ctx, username)
	if err == nil {
		fmt.Println("Usuario ya existe:", userRow.ID, userRow.Username, userRow.Email)
	} else if errors.Is(err, sql.ErrNoRows) {
		createdUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
			Username: username,
			Password: password,
			Email:    email,
		})
		if err != nil {
			log.Fatal("Error creating user:", err)
		}
		fmt.Println("Created user:", createdUser)
	} else {
		log.Fatal("Error checking existing user:", err)
	}

	fmt.Println("Server started on port", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
