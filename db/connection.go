package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() *pgxpool.Pool {
	// cargar variables .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema.")
	}

	// conexion a BD
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to constructing it if DATABASE_URL is missing
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbName,
		)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// aplicar schema.sql
	schemaBytes, err := os.ReadFile("db/schema/schema.sql")
	if err != nil {
		log.Fatalf("Error leyendo schema.sql: %v", err)
	}

	// pgx Exec supports multiple statements
	if _, err := pool.Exec(context.Background(), string(schemaBytes)); err != nil {
		log.Fatalf("Error aplicando schema: %v", err)
	}

	fmt.Println("Conexión a la base de datos exitosa (pgxpool)")
	return pool
}
