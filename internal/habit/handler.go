package habit

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
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
