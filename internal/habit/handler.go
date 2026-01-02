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

func (h *HabitHandler) HandleGetHabitList(w http.ResponseWriter, r *http.Request) {
	var (
		q         HabitListQuery
		startDate dateFilter
		endDate   dateFilter
		createdAt dateFilter
	)

	user := context.GetUser(r)
	validSort := []string{"id", "created_at", "start_date", "end_date"}

	search := utils.ReadString(r, "search", "")
	habitType := utils.ReadString(r, "type", "")
	minStartDateStr := utils.ReadString(r, "min_start", "")
	maxStartDateStr := utils.ReadString(r, "max_start", "")
	minEndDateStr := utils.ReadString(r, "min_end", "")
	maxEndtDateStr := utils.ReadString(r, "max_end", "")
	minCreatedAtStr := utils.ReadString(r, "min_created_at", "")
	maxCreatedAtStr := utils.ReadString(r, "max_created_at", "")
	privacyType := utils.ReadString(r, "status", "")

	q.sort = utils.ReadString(r, "sort", "id")
	q.pageSize = utils.ReadInt(r, "page_size", 50)
	q.page = utils.ReadInt(r, "page", 1)
	q.sortSafe = validSort
	q.userRole = user.UserRole

	err := q.validateSort()
	if err != nil {
		response.BadRequest(w, r, err, h.logger)
		return
	}

	if privacyType == "private" && user.UserRole != "admin" {
		response.Unauthorized(w, r, "You can only query public habits")
		return
	}

	q.search = search
	q.habitType = habitType
	q.privacyType = privacyType

	if minStartDateStr != "" {
		minDate, err := utils.ConvertStrToDate(minStartDateStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		startDate.minDate = minDate
	}
	if maxStartDateStr != "" {
		maxDate, err := utils.ConvertStrToDate(maxStartDateStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		startDate.maxDate = maxDate
	}

	if minEndDateStr != "" {
		minDate, err := utils.ConvertStrToDate(minEndDateStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		endDate.minDate = minDate
	}
	if maxEndtDateStr != "" {
		maxDate, err := utils.ConvertStrToDate(maxEndtDateStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		endDate.maxDate = maxDate
	}

	if minCreatedAtStr != "" {
		minDate, err := utils.ConvertStrToDate(minCreatedAtStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		createdAt.minDate = minDate
	}
	if maxCreatedAtStr != "" {
		maxDate, err := utils.ConvertStrToDate(maxCreatedAtStr)
		if err != nil {
			response.BadRequest(w, r, err, h.logger)
			return
		}
		createdAt.maxDate = maxDate
	}

	data, metaData, err := h.service.list(r.Context(), q)
	if err == errDateQuery || err == errHabitType {
		response.BadRequest(w, r, err, h.logger)
		return
	}
	if err != nil {
		response.InternalServerError(w, r, err, h.logger)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{
		"metaData": metaData,
		"data":     data,
	})
}
