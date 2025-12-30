package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	cx "github.com/NurulloMahmud/habits/pkg/context"
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

	data, err := h.service.register(r.Context(), req)
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	user, token, err := h.service.login(r.Context(), *req.Email, *req.Password)
	if err != nil {
		switch err {
		case errInvalidCredentials:
			response.BadRequest(w, r, err, h.logger)
			return
		case errUserInactive:
			response.Unauthorized(w, r, "Unauthorized")
			return
		case errUserLocked:
			response.Unauthorized(w, r, "Unauthorized")
			return
		case errMatchingPassword:
			response.InternalServerError(w, r, err, h.logger)
			return
		default:
			response.InternalServerError(w, r, err, h.logger)
			return
		}
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{
		"access_token": token,
		"user":         user,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	err = req.validateUpdateUserRequest()
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	user := cx.GetUser(r)
	err = h.service.update(r.Context(), user.ID, req)
	
	if err != nil {
		switch err {
		case errInvalidCredentials:
			response.BadRequest(w, r, err, h.logger)
			return
		case errEmailTaken:
			response.BadRequest(w, r, err, h.logger)
			return
		default:
			response.InternalServerError(w, r, err, h.logger)
			return
		}
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"message": "user updated successfully"})
}
