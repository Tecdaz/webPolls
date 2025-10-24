package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"webpolls/services"
	"webpolls/utils"
)

// PollHandler ahora depende de PollService
type PollHandler struct {
	service *services.PollService
}

// NewPollHandler ahora inyecta PollService
func NewPollHandler(service *services.PollService) *PollHandler {
	return &PollHandler{service: service}
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req services.PollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Cuerpo json invalido")
		return
	}

	data, err := h.service.CreatePoll(r.Context(), req)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondWithData(w, http.StatusCreated, data, "Encuesta creada correctamente")
}

func (h *PollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ConvertTo32(idStr)
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
	id, err := utils.ConvertTo32(idStr)
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

// modificado para traer las opciones junto con la encuesta
func (h *PollHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.service.GetPolls(r.Context())
	if err != nil {
		log.Printf("Error getting polls: %v", err) // Agregar log
		RespondWithError(w, http.StatusInternalServerError, "No se pudieron obtener las encuestas")
		return
	}

	log.Printf("Polls retrieved: %+v", polls) // Agregar log para debug
	RespondWithData(w, http.StatusOK, polls, "Encuestas obtenidas correctamente")
}

func (h *PollHandler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ConvertTo32(idStr)
	if err != nil {
		log.Printf("Error converting option ID: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Id de opción inválido")
		return
	}

	var req services.OptionResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		RespondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	log.Printf("Updating option %d with content=%v", id, req.Content)

	data, err := h.service.UpdateOption(r.Context(), req)
	if err != nil {
		log.Printf("Error updating option: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error al actualizar opción")
		return
	}

	log.Printf("Option %d updated successfully to content=%v", id, req.Content)
	RespondWithData(w, http.StatusOK, data, "Opción actualizada correctamente")
}

func (h *PollHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	poll_idStr := r.PathValue("poll_id")
	id, err := utils.ConvertTo32(idStr)
	if err != nil {
		log.Printf("Error converting option ID: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Id de opción inválido")
		return
	}
	poll_id, err := utils.ConvertTo32(poll_idStr)
	if err != nil {
		log.Printf("Error converting poll ID: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Id de encuesta inválido")
		return
	}

	err = h.service.DeleteOption(r.Context(), id, poll_id)
	if err != nil {
		log.Printf("Error deleting option: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error al eliminar opción")
		return
	}

	log.Printf("Option %d deleted successfully", id)
	RespondWithData(w, http.StatusOK, nil, "Opción eliminada correctamente")
}
