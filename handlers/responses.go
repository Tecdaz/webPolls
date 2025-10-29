package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ApiResponse es la estructura estándar para todas las respuestas de la API.
// asi el frontend conoce de antemano la estructura que recibe por cada respuesta para siempre buscar los datos en el campo data y los mensajes de error en el campo error.
type ApiResponse struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// RespondWithError envía una respuesta de error JSON estandarizada.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// Cuando hay un error, Data es nil.
	payload := ApiResponse{Error: message, Data: nil}
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error al codificar respuesta JSON: %v", err)
	}
}

// RespondWithData envía una respuesta JSON exitosa estandarizada.
func RespondWithData(w http.ResponseWriter, code int, dataPayload interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// Cuando la respuesta es exitosa, Error está vacío y se omite.
	payload := ApiResponse{Data: dataPayload, Message: message}
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error al codificar respuesta JSON: %v", err)
	}
}
