package handlers

import (
	"encoding/json"
	"net/http"

	"webpolls/components"
	"webpolls/services"
	"webpolls/utils"
	"webpolls/views"

	"github.com/jackc/pgx/v5"
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
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	// Redirect to login page
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
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
		if err == pgx.ErrNoRows {
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

func (h *userHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	// Si ya está logueado, redirigir al home
	if utils.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := views.AuthLayout(views.Login(), "Iniciar Sesión - Webpolls").Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *userHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		components.Toast("Invalid form data", true).Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.service.Authenticate(r.Context(), email, password)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusUnauthorized)
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	// Crear sesión
	session := utils.GetSession(r)
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.Id
	session.Values["username"] = user.Username
	utils.SaveSession(w, r, session)

	// Redirigir al home usando HTMX
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session := utils.GetSession(r)
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	utils.SaveSession(w, r, session)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *userHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	// Si ya está logueado, redirigir al home
	if utils.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := views.AuthLayout(views.Register(), "Registrarse - Webpolls").Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
