package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"webpolls/dataconvertion"
	"webpolls/services"
)

// PollHandler ahora depende de PollService
type PollHandler struct {
	service *services.PollService
}

// NewPollHandler ahora inyecta PollService
func NewPollHandler(service *services.PollService) *PollHandler {
	return &PollHandler{service: service}
}

// CreatePollRequest sigue siendo relevante para el handler, para decodificar el JSON
type CreatePollRequest struct {
	Question string                      `json:"question"`
	Options  []services.OptionCreatePoll `json:"options"`
	UserID   int32                       `json:"user_id"`
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Cuerpo json invalido")
		return
	}

	params := services.CreatePollParams{
		Question: req.Question,
		Options:  req.Options,
		UserID:   req.UserID,
	}

	data, err := h.service.CreatePoll(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondWithData(w, http.StatusCreated, data, "Encuesta creada correctamente")
}

func (h *PollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := dataconvertion.ConvertTo32(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de poll invalido")
		return
	}

	err = h.service.DeletePoll(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithData(w, http.StatusOK, nil, "Encuesta eliminada correctamente")
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := dataconvertion.ConvertTo32(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de encuesta invalido")
		return
	}

	poll, err := h.service.GetPollByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Encuesta no encontrada")
		} else {
			log.Println("DB error:", err)
			RespondWithError(w, http.StatusInternalServerError, "Error al obtener encuesta")
		}
		return
	}

	RespondWithData(w, http.StatusOK, poll, "Encuesta obtenida correctamente")
}

func (h *PollHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.service.GetPolls(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "No se pudieron obtener las encuestas")
		return
	}

	RespondWithData(w, http.StatusOK, polls, "Encuestas obtenidas correctamente")
}
