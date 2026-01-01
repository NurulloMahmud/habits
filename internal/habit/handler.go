package habit

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
	"github.com/NurulloMahmud/habits/pkg/utils"
)

type HabitHandler struct {
	service Service
	logger  *log.Logger
}

func NewHandler(s Service, log *log.Logger) *HabitHandler {
	return &HabitHandler{
		service: s,
		logger:  log,
	}
}

func (h *HabitHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req createHabitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	user := context.GetUser(r)
	err = req.validateCreateRequest(user.ID)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	data, err := h.service.create(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.Envelope{"data": data})
}

func (h *HabitHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var req updateHabitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}
	req.ID = int(habitID)

	user := context.GetUser(r)
	data, err := h.service.update(r.Context(), user.ID, req)

	if err != nil {
		switch err {
		case errNoHabitFound:
			response.BadRequest(w, r, err, h.logger)
			return
		case errNotOwner:
			response.Forbidden(w, r, err.Error())
			return
		case errTypeChange:
			response.BadRequest(w, r, err, h.logger)
			return
		case errNameEmpty:
			response.BadRequest(w, r, err, h.logger)
			return
		case errDescEmpty:
			response.BadRequest(w, r, err, h.logger)
			return
		case errInvalidDates:
			response.BadRequest(w, r, err, h.logger)
			return
		case errInvalidStatus:
			response.BadRequest(w, r, err, h.logger)
			return
		default:
			response.InternalServerError(w, r, err, h.logger)
			return
		}
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"data": data, "message": "habit updated successfully"})
}

func (h *HabitHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	user := context.GetUser(r)
	err = h.service.delete(r.Context(), *user, habitID)
	if err != nil {
		if errors.Is(err, errNoHabitFound) {
			response.BadRequest(w, r, err, h.logger)
			return
		} else if errors.Is(err, errNotOwner) {
			response.Forbidden(w, r, "You are not the owner of this habit")
			return
		}
		response.InternalServerError(w, r, err, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusNoContent, response.Envelope{"message": "habit deleted successfully"})
}

func (h *HabitHandler) HandleGetPrivateHabit(w http.ResponseWriter, r *http.Request) {
	identifier, err := utils.ReadIdentifierParam(r)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	habit, err := h.service.repo.get(r.Context(), 0, identifier)
	if err != nil {
		response.InternalServerError(w, r, err, h.logger)
		return
	}
	if habit == nil {
		response.BadRequest(w, r, errNoHabitFound, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"data": habit})
}
