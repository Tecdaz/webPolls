package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"webpolls/services"
	"webpolls/utils"
)

// userHandler ahora depende de UserService
type userHandler struct {
	service *services.UserService
}

// NewUserHandler ahora inyecta UserService
func NewUserHandler(service *services.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req services.UserRequest
	contentType := r.Header.Get("Content-Type")
	// modo json para peticiones de APIs y HTMX
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
			return
		}
	} else { //envío tradicional de forms HTML
		if err := r.ParseForm(); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid form data")
			return
		}
		req = services.UserRequest{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error()) 
		return
	}

	RespondWithData(w, http.StatusCreated, user, "Usuario creado correctamente")
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := utils.ConvertTo32(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de usuario invalido")
		return
	}

	deletedUsername, err := h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//API
	RespondWithData(w, http.StatusOK, map[string]string{"username": deletedUsername}, "Usuario eliminado correctamente")
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := utils.ConvertTo32(id)
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

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := utils.ConvertTo32(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de usuario invalido")
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Payload de actualizacion invalido")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), userID, req)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	RespondWithData(w, http.StatusOK, user, "Usuario actualizado correctamente")
}

func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "No se pudieron obtener los usuarios")
		return
	}

	RespondWithData(w, http.StatusOK, users, "Usuarios obtenidos correctamente")
}
