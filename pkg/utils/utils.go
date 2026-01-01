package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}

func ReadIdentifierParam(r *http.Request) (string, error) {
	idParam := chi.URLParam(r, "identifier")
	if idParam == "" {
		return "", errors.New("invalid id parameter")
	}

	return idParam, nil
}

func ReadString(r *http.Request, key string, defaultValue string) string {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func ReadBool(r *http.Request, key string) (*bool, error) {
	str := r.URL.Query().Get(key)

	if str == "" {
		return nil, nil
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil, errors.New("invalid boolean field provided in query param")
	}
	return &b, nil
}

func ReadCSV(r *http.Request, key string, defaultValue []string) []string {
	csv := r.URL.Query().Get(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func ReadInt(r *http.Request, key string, defaultValue int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
