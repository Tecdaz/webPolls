package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
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
	return db
}
