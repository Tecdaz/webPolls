package handlers

import (
	"net/http"
	"webpolls/services"
	"webpolls/views"
)

type homeHandler struct {
	service *services.UserService
}

func NewHomeHandler(service *services.UserService) *homeHandler {
	return &homeHandler{service: service}
}

//maneja las solicitudes a la pagina principal, si se accede por primera vez se sirve el HTML completo, si la solicitud proviene de la navegacion interna se renderiza y devuelve el contenido especifico
func (h *homeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		err := views.Home().Render(r.Context(), w)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	err := views.Index(views.Home(), "Webpolls").Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
