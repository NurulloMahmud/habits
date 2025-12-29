package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error, log *log.Logger) {
	requestURL := r.URL.String()
	requestMethod := r.Method
	log.Printf("[ERROR] SERVER ERROR!\nURL: %s\nMethod: %s\nError: %v\n", requestURL, requestMethod, err.Error())
	WriteJSON(w, http.StatusInternalServerError, Envelope{"error": "internal server error"})
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, log *log.Logger) {
	requestURL := r.URL.String()
	requestMethod := r.Method
	log.Printf("[ERROR] BAD REQUEST!\nURL: %s\nMethod: %s\nError: %v\n", requestURL, requestMethod, err.Error())
	WriteJSON(w, http.StatusBadRequest, Envelope{"error": err.Error()})
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	requestURL := r.URL.String()
	requestMethod := r.Method
	log.Printf("[ERROR] UNAUTHORIZED REQUEST!\nURL: %s\nMethod: %s\nError: %v\n", requestURL, requestMethod, msg)
	WriteJSON(w, http.StatusUnauthorized, Envelope{"error": msg})
}

func Forbidden(w http.ResponseWriter, r *http.Request, msg string) {
	requestURL := r.URL.String()
	requestMethod := r.Method
	log.Printf("[ERROR] FORBIDDEN REQUEST!\nURL: %s\nMethod: %s\nError: %v\n", requestURL, requestMethod, msg)
	WriteJSON(w, http.StatusForbidden, Envelope{"error": msg})
}
