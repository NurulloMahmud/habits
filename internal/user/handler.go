package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/response"
)

var (
	internalServerError = errors.New("internal server error")
)

type UserHandler struct {
	service UserService
	logger  *log.Logger
}

func NewHandler(s UserService, log *log.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  log,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	err = req.validateRegister()
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	data, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, errEmailTaken) {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		response.InternalServerError(w, r, err, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.Envelope{"data": data})
}
