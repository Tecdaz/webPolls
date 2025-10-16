package handlers

import (
	"database/sql"
	"log"
	sqlc "webpolls/db/sqlc"

	"webpolls/dataconvertion"

	"github.com/gin-gonic/gin"
)

// para inicializar el handler desde el main
type pollHandler struct {
	queries *sqlc.Queries
}

func NewPollHandler(queries *sqlc.Queries) *pollHandler {
	return &pollHandler{queries: queries}
}

// objeto que recibe el json desde la solicitud
type CreatePollRequest struct {
	Question string   `json:"question" binding:"required"`
	Options  []string `json:"options" binding:"required,min=2"`
}

func (h *pollHandler) CreatePoll(c *gin.Context) {
	var req CreatePollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//crea la encuesta
	poll, err := h.queries.CreatePoll(c, sqlc.CreatePollParams{
		Title:  req.Question,
		UserID: 1, // 1 por que todavia no hice la parte de los usuarios
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create poll"})
		return
	}

	/*
		en esta parte creo las opciones, porque las opciones dependen de la encuesta
		no tiene sentido hacer un handler de opciones suelto
	*/

	var options []sqlc.CreateOptionRow
	for _, optionContent := range req.Options {
		option, err := h.queries.CreateOption(c, sqlc.CreateOptionParams{
			Content: optionContent,
			Correct: sql.NullBool{Bool: false, Valid: true},
			PollID:  poll.ID,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create option"})
			return
		}
		options = append(options, option)
	}

	c.JSON(201, gin.H{
		"poll":    poll,
		"options": options,
	})
}

func (h *pollHandler) DeletePoll(c *gin.Context) {
	pollID := c.Param("id")
	pollIDInt32, err := dataconvertion.ConvertTo32(pollID)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid poll ID"})
		return
	}

	err = h.queries.DeletePoll(c, pollIDInt32)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete poll"})
		return
	}
	c.JSON(200, gin.H{"message": "Poll deleted successfully"})
}

func (h *pollHandler) GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	pollIDInt32, err := dataconvertion.ConvertTo32(pollID)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid poll ID"})
		return
	}
	poll, err := h.queries.GetPollByID(c, pollIDInt32)
	if err != nil {
		log.Println("DB error:", err)
		c.JSON(500, gin.H{"error": "Failed to get poll"})
		return
	}
	c.JSON(200, gin.H{"poll": poll})
}
