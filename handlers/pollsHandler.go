package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"webpolls/components"
	"webpolls/middleware"
	"webpolls/services"
	"webpolls/utils"
	"webpolls/views"
)

// PollHandler ahora depende de PollService
type PollHandler struct {
	service *services.PollService
	sse     *services.SSEBroker
}

// NewPollHandler ahora inyecta PollService y SSEBroker
func NewPollHandler(service *services.PollService, sse *services.SSEBroker) *PollHandler {
	return &PollHandler{service: service, sse: sse}
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req services.PollRequest

	if err := r.ParseForm(); err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		components.Toast("Cuerpo forma invalido", true).Render(r.Context(), w)
		return
	}
	var options []services.OptionRequest
	formOptions := r.Form["options"]
	for _, opt := range formOptions {
		options = append(options, services.OptionRequest{Content: opt})
	}
	userId := r.Context().Value(middleware.UserIDKey).(int32)

	req = services.PollRequest{
		Question: r.FormValue("question"),
		UserID:   userId,
		Options:  options,
	}

	_, err := h.service.CreatePoll(r.Context(), req)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	//se llama a esto para traer todas las polls y mandarlas al renderizado
	polls, err := h.service.GetPollsByUser(r.Context(), userId, userId)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}

	err = views.PollList(polls, true).Render(r.Context(), w)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		components.Toast(err.Error(), true).Render(r.Context(), w)
		return
	}
	components.Toast("Encuesta creada correctamente", false).Render(r.Context(), w)
}

func (h *PollHandler) GetMyPolls(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int32)

	polls, err := h.service.GetPollsByUser(r.Context(), userId, userId)
	if err != nil {
		log.Printf("Error getting user polls: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "No se pudieron obtener las encuestas")
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		err := views.MyPolls(polls).Render(r.Context(), w)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	err = views.Layout(views.MyPolls(polls), "Mis Polls - Webpolls", utils.IsAuthenticated(r)).Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
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

	//WEB
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	//API
	RespondWithData(w, http.StatusOK, nil, "Encuesta eliminada correctamente")
}

func (h *PollHandler) GetPollPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ConvertTo32(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de encuesta invalido")
		return
	}

	var userID *int32
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		uid := val.(int32)
		userID = &uid
	}

	poll, err := h.service.GetPollByID(r.Context(), id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Encuesta no encontrada")
		} else {
			log.Println("DB error:", err)
			RespondWithError(w, http.StatusInternalServerError, "Error al obtener encuesta")
		}
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		views.PollDetailContent(poll, userID != nil).Render(r.Context(), w)
		return
	}

	err = views.Layout(views.PollDetail(poll, userID != nil), "Encuesta - Webpolls", userID != nil).Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *PollHandler) Vote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	pollID, err := utils.ConvertTo32(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de encuesta invalido")
		return
	}

	if err := r.ParseForm(); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Formulario inválido")
		return
	}
	optionIDStr := r.FormValue("option_id")
	optionID, err := utils.ConvertTo32(optionIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Id de opción inválido")
		return
	}

	var userID int32
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(int32)
	} else {
		RespondWithError(w, http.StatusUnauthorized, "Debes iniciar sesión para votar")
		return
	}

	err = h.service.Vote(r.Context(), pollID, optionID, userID)
	if err != nil {
		log.Printf("Error voting: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error al registrar voto")
		return
	}

	// Broadcast update via SSE
	// We need to render the poll detail content for ALL users, but wait...
	// SSE sends the SAME data to everyone. But the view depends on "UserVotedOptionID".
	// If we broadcast the full HTML, it will be specific to a user (or generic).
	// Strategy: Broadcast an event "poll_updated" with the poll ID.
	// The client receives it and triggers an HX-GET to refresh the content.
	// OR: Broadcast the "Results" part which is common (percentages).
	// But the "Your Vote" badge is specific.
	//
	// Let's use the "trigger event" strategy.
	// We send a simple event: "poll_updated:<id>"
	// The client listens and reloads.
	//
	// Wait, the requirement says "tiempo real". If we trigger a reload, it's a new request.
	// That's fine for HTMX.
	//
	// Let's try to send the updated HTML for the *results* only?
	// The results (bars/percentages) are the same for everyone.
	// The "Your Vote" badge is local.
	//
	// If we send HTML, we can't easily personalize it for each listener.
	// So, we should send an event that triggers a refresh on the client.
	// HTMX SSE extension supports this: `hx-trigger="sse:message"`.
	//
	// Let's send a custom event name: "update_poll_<id>"
	//
	// In the view: <div hx-ext="sse" sse-connect="/events" sse-swap="update_poll_<id>" hx-get="/polls/<id>" hx-trigger="sse:update_poll_<id>">
	//
	// Actually, `sse-swap` expects the event data to be HTML to swap.
	// `hx-trigger` allows triggering a standard HTMX request.
	//
	// Let's go with `hx-trigger`.
	// Event name: "poll_update"
	// Data: {"poll_id": 123} (JSON)
	//
	// But standard HTMX SSE trigger works on event names.
	// Server sends:
	// event: poll_update_123
	// data: {}
	//
	// Client: hx-trigger="sse:poll_update_123"
	//
	// Broadcast update via SSE
	// Fetch updated poll data for broadcast (generic, no user specific data needed for stats)
	updatedPoll, err := h.service.GetPollByID(r.Context(), pollID, nil)
	if err == nil {
		// Create a simple struct for the payload
		type PollStats struct {
			TotalVotes int64 `json:"total_votes"`
			Options    []struct {
				ID         int32   `json:"id"`
				VoteCount  int64   `json:"vote_count"`
				Percentage float64 `json:"percentage"`
			} `json:"options"`
		}

		stats := PollStats{
			TotalVotes: updatedPoll.TotalVotes,
		}
		for _, opt := range updatedPoll.Options {
			stats.Options = append(stats.Options, struct {
				ID         int32   `json:"id"`
				VoteCount  int64   `json:"vote_count"`
				Percentage float64 `json:"percentage"`
			}{
				ID:         opt.ID,
				VoteCount:  opt.VoteCount,
				Percentage: opt.Percentage,
			})
		}

		jsonData, _ := json.Marshal(stats)
		h.sse.Broadcast([]byte(fmt.Sprintf("event: poll_update_%d\ndata: %s\n\n", pollID, jsonData)))
	} else {
		// Fallback if fetch fails
		h.sse.Broadcast([]byte(fmt.Sprintf("event: poll_update_%d\ndata: {}\n\n", pollID)))
	}

	poll, err := h.service.GetPollByID(r.Context(), pollID, &userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error al obtener datos actualizados")
		return
	}

	views.PollDetailContent(poll, true).Render(r.Context(), w)
}

func (h *PollHandler) SSE(w http.ResponseWriter, r *http.Request) {
	h.sse.ServeHTTP(w, r)
}

// modificado para traer las opciones junto con la encuesta
func (h *PollHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	var userID int32
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(int32)
	}

	polls, err := h.service.GetPolls(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting polls: %v", err) // Agregar log
		RespondWithError(w, http.StatusInternalServerError, "No se pudieron obtener las encuestas")
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		err := views.Polls(polls).Render(r.Context(), w)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	err = views.Layout(views.Polls(polls), "Webpolls - Polls", utils.IsAuthenticated(r)).Render(r.Context(), w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
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

func (h *PollHandler) GetPollOptionInput(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count, _ := utils.ConvertTo32(countStr)

	if count >= 4 {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		components.Toast("Máximo 4 opciones permitidas", true).Render(r.Context(), w)
		return
	}

	views.PollOptionInput().Render(r.Context(), w)
}
