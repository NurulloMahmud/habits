package habitmember

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/response"
)

type Handler struct {
	service Service
	logger  *log.Logger
}

func NewHandler(s Service, log *log.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  log,
	}
}

func (h *Handler) HandleJoinHabit(w http.ResponseWriter, r *http.Request) {
	var req habitMemberCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	msg, err := h.service.joinHabit(r.Context(), req)
	if err != nil {
		if err == sql.ErrNoRows || err == errAlreadyMember {
			response.BadRequest(w, r, err, h.logger)
			return
		}

		response.InternalServerError(w, r, err, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.Envelope{"message": msg})
}
