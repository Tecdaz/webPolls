package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"webpolls.com/webpolls/models"
)

type UserHandler struct {
	DB *gorm.DB
}

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

func connectToDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=2002 dbname=webpoll port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Migrar esquema con los modelos corregidos
	err = db.AutoMigrate(&models.User{}, &models.Poll{}, &models.Option{})
	if err != nil {
		log.Fatal("Error migrating database schema:", err)
	}

	fmt.Println("Database connection successful")
	return db
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveIndex)

	port := ":8080"
	db := connectToDB()
	// Inicializar handler con la conexión DB
	userHandler := &UserHandler{DB: db}

	http.HandleFunc("/api/users", userHandler.CreateUser)

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// Crear usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
