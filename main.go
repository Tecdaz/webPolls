package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"context"

	sqlc "webpolls.com/webpolls/db/sqlc"

	_ "github.com/lib/pq"
)

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

	//conexion a la base de datos
	constStr := "user=joaquin password=1999 dbname=webpoll host=localhost port=5432 sslmode=disable"
	db, errDb := sql.Open("postgres", constStr)

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

	fmt.Println("Successfully connected to database!")

	createdUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: "testuser",
		Password: "password123",
		Email:    "jklsajdlaskj@skd.com",
	})

	if err != nil {
		log.Fatal("Error creating user:", err)
	}

	fmt.Println("Created user:", createdUser)

	fmt.Println("Server started on port", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
