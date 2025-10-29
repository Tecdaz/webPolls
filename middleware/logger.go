package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

//es un wrapper que captura el código de estado y el cuerpo de la respuesta.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer // Campo para capturar el cuerpo de la respuesta
}

//crea una nueva instancia de nuestro wrapper.
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer), 
	}
}

// captura el código de estado.
func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// captura el cuerpo de la respuesta y lo escribe en la respuesta original.
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	w.body.Write(b) // guardamos los bytes en nuestro buffer
	return w.ResponseWriter.Write(b) // pasamos los bytes al ResponseWriter original
}

//es el middleware que registra la información de cada petición.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		//creamos una copia del ResponseWriter original que nos permite ver y guardar lo que se va a enviar al cliente sin afectar la respuesta real
		wrapper := newResponseWriterWrapper(w)

		//lectura del cuerpo de la petición
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		// restaurmaos el cuerpo para que los siguientes handlers puedan leerlo
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// registramos la petición entrante
		log.Printf("--> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		if len(bodyBytes) > 0 {
			log.Printf("    Request Body: %s", bodyBytes)
		}

		next.ServeHTTP(wrapper, r) //siguiente handler en la cadena con nuestro wrapper

		// registramos la respuesta saliente
		duration := time.Since(start)
		log.Printf("<-- %d %s in %s", wrapper.statusCode, http.StatusText(wrapper.statusCode), duration)

		if wrapper.body.Len() > 0 {
			log.Printf("    Response Body: %s", wrapper.body.String())
		}
	})
}
