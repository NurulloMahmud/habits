package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/context"
	cx "github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
	"github.com/NurulloMahmud/habits/pkg/utils"
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

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	var input ListUserInput

	user := context.GetUser(r)
	if user.UserRole != "admin" {
		response.Unauthorized(w, r, "Unauthorized")
		return
	}

	input.Search = utils.ReadString(r, "search", "")
	input.Filter.Page = utils.ReadInt(r, "page", 1)
	input.Filter.PageSize = utils.ReadInt(r, "page_size", 50)
	input.Filter.Sort = utils.ReadString(r, "sort", "id")
	input.Filter.SortSafeList = []string{"id", "email", "first_name", "last_name"}

	is_active, err := utils.ReadBool(r, "is_active")
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}
	is_locked, err := utils.ReadBool(r, "is_locked")
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	input.IsActive = is_active
	input.IsLocked = is_locked

	err = input.ValidateInput()
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	users, metadata, err := h.service.list(r.Context(), input)
	if err != nil {
		response.InternalServerError(w, r, err, h.logger)
		return
	}

	if users == nil {
		response.WriteJSON(w, http.StatusOK, response.Envelope{"result": []any{}, "message": "no user found"})
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"result": users, "metadata": metadata})
}
