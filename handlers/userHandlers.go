package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"webpolls/dataconvertion"
	sqlc "webpolls/db/sqlc"
	"webpolls/services"
)

// userHandler ahora depende de UserService
type userHandler struct {
	service *services.UserService
}

// NewUserHandler ahora inyecta UserService
func NewUserHandler(service *services.UserService) *userHandler {
	return &userHandler{service: service}
}

// CreateUserRequest sigue siendo relevante para el handler, para decodificar el JSON
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Payload de creacion invalido")
		return
	}

	params := sqlc.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// La lógica de negocio ahora se llama desde el servicio
	user, err := h.service.CreateUser(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error()) // El servicio devuelve errores de negocio claros
		return
	}

	RespondWithData(w, http.StatusCreated, user, "Usuario creado correctamente")
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := dataconvertion.ConvertTo32(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de usuario invalido")
		return
	}

	deletedUsername, err := h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithData(w, http.StatusOK, map[string]string{"username": deletedUsername}, "Usuario eliminado correctamente")
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := dataconvertion.ConvertTo32(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de usuario invalido")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Usuario no encontrado")
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithData(w, http.StatusOK, user, "Usuario obtenido correctamente")
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := dataconvertion.ConvertTo32(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de usuario invalido")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Payload de actualizacion invalido")
		return
	}

	// Convertimos el modelo del handler al modelo del servicio
	serviceParams := services.UpdateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.service.UpdateUser(r.Context(), userID, serviceParams)
	if err != nil {
		// El servicio devuelve errores de negocio que podemos mapear a códigos de estado HTTP
		if err.Error() == "usuario no encontrado" {
			RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	RespondWithData(w, http.StatusOK, user, "Usuario actualizado correctamente")
}
