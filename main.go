package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"webpolls/handlers"

	sqlc "webpolls/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar variables de entorno
	_ = godotenv.Load()

	// Configurar conexión a la DB
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable",
		dbUser, dbPassword, dbName, dbHost)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Ejecutar esquema
	schemaBytes, err := os.ReadFile("db/schema/schema.sql")
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}
	if _, err := db.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Error applying schema: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	queries := sqlc.New(db)

	//iniciar handlers
	pollHandler := handlers.NewPollHandler(queries)
	userHandler := handlers.NewUserHandler(queries)

	// --- Gin Setup ---
	router := gin.Default()

	// Servir archivos estáticos
	router.Static("/static", "./static")

	// Ruta principal
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	//llamados
	router.POST("/polls/create", pollHandler.CreatePoll)
	router.DELETE("/polls/:id", pollHandler.DeletePoll)
	router.GET("/polls/:id", pollHandler.GetPoll)
	router.POST("/users/create", userHandler.CreateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("users/:id", userHandler.UpdateUser)

	// Inicializar servidor
	port := ":8080"
	fmt.Println("Server started on port", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
