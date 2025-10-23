package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"webpolls/dataconvertion"
	db "webpolls/db/sqlc"
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

// agregado para incluir las opciones completas
type Option struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Correct bool   `json:"correct"`
}

type PollWithOptions struct {
	PollID  int      `json:"poll_id"`
	Title   string   `json:"title"`
	UserID  int32    `json:"user_id"`
	Options []Option `json:"options"`
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

// modificado para traer las opciones junto con la encuesta
func (h *PollHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.service.GetPollsWithOptions(r.Context())
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
	id, err := dataconvertion.ConvertTo32(idStr)
	if err != nil {
		log.Printf("Error converting option ID: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Id de opción inválido")
		return
	}

	var body struct {
		Correct bool `json:"correct"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		RespondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	log.Printf("Updating option %d with correct=%v", id, body.Correct)

	// Traigo la opción actual para obtener el poll_id
	opt, err := h.service.Queries.GetOptionByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting option by ID: %v", err)
		RespondWithError(w, http.StatusNotFound, "Opción no encontrada")
		return
	}

	// Si estamos marcando como correcto, primero desmarcamos todas las otras opciones de la misma encuesta
	if body.Correct {
		err = h.service.Queries.UnsetOtherOptionsCorrect(r.Context(), db.UnsetOtherOptionsCorrectParams{
			PollID: opt.PollID,
			ID:     id,
		})
		if err != nil {
			log.Printf("Error unsetting other options: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Error al actualizar otras opciones")
			return
		}
	}

	// Actualizo la opción específica
	err = h.service.Queries.UpdateOption(r.Context(), db.UpdateOptionParams{
		ID:      id,
		Content: opt.Content,
		Correct: sql.NullBool{Bool: body.Correct, Valid: true},
	})
	if err != nil {
		log.Printf("Error updating option: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error al actualizar opción")
		return
	}

	log.Printf("Option %d updated successfully to correct=%v", id, body.Correct)
	RespondWithData(w, http.StatusOK, nil, "Opción actualizada correctamente")
}
