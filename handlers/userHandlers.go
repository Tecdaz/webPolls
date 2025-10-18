package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	sqlc "webpolls/db/sqlc"
	"webpolls/dataconvertion"
)

type userHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(queries *sqlc.Queries) *userHandler {
	return &userHandler{queries: queries}
}

//create user

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //solo uso de POST
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest //verificacion de cuerpo json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// validacion de campos
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// asegurar que no exista ni el usuario ni el mail
	_, err := h.queries.GetUserByUsername(r.Context(), req.Username) //uso de r.context para que las consultas a la bdd se cancelen automaticamente si el cliente corta la conexion
	if err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	user, err := h.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user": user})
}

// delete user
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := dataconvertion.ConvertTo32(id) //asi lo espera sqlc
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteUser(r.Context(), userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// get user
func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := dataconvertion.ConvertTo32(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user": user})
}

/*
DISCLAIMER: revisar esta parte por que ya estaba re quemado y la hice con gpt
*/

// update user
type UpdateUserRequest struct {
	Username *string `json:"username"` //campo puntero para distinguir entre ausente y vacio
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := dataconvertion.ConvertTo32(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Verificar que exista
	existingUser, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var username, email, password sql.NullString //para manejar campos opcionales

	if req.Username != nil {
		userByUsername, err := h.queries.GetUserByUsername(r.Context(), *req.Username)
		// Si se encontró un usuario con ese username y NO es el mismo usuario -> error
		if err == nil && userByUsername.ID != existingUser.ID {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
		username = sql.NullString{String: *req.Username, Valid: true}
	}

	if req.Email != nil {
		userByEmail, err := h.queries.GetUserByEmail(r.Context(), *req.Email)
		// Si se encontró un usuario con ese email y NO es el mismo usuario -> error
		if err == nil && userByEmail.ID != existingUser.ID {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}
		email = sql.NullString{String: *req.Email, Valid: true}
	}

	if req.Password != nil {
		password = sql.NullString{String: *req.Password, Valid: true}
	}

	user, err := h.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		ID:      userID,
		Column2: username,
		Column3: email,
		Column4: password,
	})
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user": user})
}