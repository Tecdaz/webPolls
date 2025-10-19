package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// responseWriterWrapper es un wrapper que captura el código de estado y el cuerpo de la respuesta.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer // Campo para capturar el cuerpo de la respuesta
}

// newResponseWriterWrapper crea una nueva instancia de nuestro wrapper.
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer), // Inicializamos el buffer
	}
}

// WriteHeader captura el código de estado.
func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captura el cuerpo de la respuesta y lo escribe en la respuesta original.
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	w.body.Write(b) // Guardamos los bytes en nuestro buffer
	return w.ResponseWriter.Write(b) // Pasamos los bytes al ResponseWriter original
}

// LoggingMiddleware es el middleware que registra la información de cada petición.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Envolvemos el ResponseWriter original para poder capturar estado y cuerpo
		wrapper := newResponseWriterWrapper(w)

		// Leemos el cuerpo de la petición para poder registrarlo
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		// Restauramos el cuerpo para que los siguientes handlers puedan leerlo
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Registramos la petición entrante
		log.Printf("--> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		// Solo registramos el cuerpo si existe
		if len(bodyBytes) > 0 {
			log.Printf("    Request Body: %s", bodyBytes)
		}

		// Llamamos al siguiente handler en la cadena con nuestro wrapper
		next.ServeHTTP(wrapper, r)

		// Registramos la respuesta saliente
		duration := time.Since(start)
		log.Printf("<-- %d %s in %s", wrapper.statusCode, http.StatusText(wrapper.statusCode), duration)

		// Registramos el cuerpo de la respuesta si no está vacío
		if wrapper.body.Len() > 0 {
			log.Printf("    Response Body: %s", wrapper.body.String())
		}
	})
}
