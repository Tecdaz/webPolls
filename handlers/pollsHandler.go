package handlers

import (
	"database/sql"
	"log"
	sqlc "webpolls/db/sqlc"

	"webpolls/dataconvertion"

	"encoding/json"
	"net/http"
	"strconv"
)

// para inicializar el handler desde el main
type PollHandler struct {
	queries *sqlc.Queries
}

func NewPollHandler(queries *sqlc.Queries) *PollHandler {
	return &PollHandler{queries: queries}
}

// objeto que recibe el json desde la solicitud
type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest

	if req.Question == "" || len(req.Options) < 2 {
    	http.Error(w, "Missing question or not enough options", http.StatusBadRequest)
    	return
	}

	//crea la encuesta
	poll, err := h.queries.CreatePoll(r.Context(), sqlc.CreatePollParams{
		Title:  req.Question,
		UserID: 1, // 1 por que todavia no hice la parte de los usuarios
	})

	if err != nil {
		http.Error(w, "Failed to create poll", http.StatusInternalServerError)
		return
	}

	/*
		en esta parte creo las opciones, porque las opciones dependen de la encuesta
		no tiene sentido hacer un handler de opciones suelto
	*/

	var options []sqlc.CreateOptionRow
	for _, optionContent := range req.Options {
		option, err := h.queries.CreateOption(r.Context(), sqlc.CreateOptionParams{
			Content: optionContent,
			Correct: sql.NullBool{Bool: false, Valid: true},
			PollID:  poll.ID,
		})
		if err != nil {
			http.Error(w, "Failed to create option", http.StatusInternalServerError)
			return
		}
		options = append(options, option)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
    	"poll":    poll,
    	"options": options,
	})

}

func (h *PollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := dataconvertion.ConvertTo32(idStr)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeletePoll(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete poll", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Poll deleted successfully",
	})
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	poll, err := h.queries.GetPollByID(r.Context(), int32(id))
	if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "Failed to get poll", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"poll": poll,
	})
}