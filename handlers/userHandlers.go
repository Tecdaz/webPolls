package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"webpolls/components"
	"webpolls/services"
	"webpolls/utils"
	"webpolls/views"
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

	if err := r.ParseForm(); err != nil {
		w.Header().Set("HX-Reswap", "none")
		components.Toast("Invalid form data", true).Render(r.Context(), w)
		return
	}
	req = services.UserRequest{
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	_, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	//nuevamente llamo a esto para traer los usuarios y renderizarlos
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	err = views.UserList(users).Render(r.Context(), w)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}
	components.Toast("Usuario creado correctamente", false).Render(r.Context(), w)
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

	//WEB
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
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

	//WEB
	if r.Header.Get("HX-Request") == "true" {
		err = views.Users(users).Render(r.Context(), w)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	err = views.Layout(views.Users(users), "Usuarios - Webpolls").Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
