package handler

import (
	"errors"
	"net/http"
	"strconv"

	"golang-test/internal/dto"
	"golang-test/internal/middleware"
	"golang-test/internal/service"
	"golang-test/pkg/response"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	calendarService service.CalendarService
}

func NewCalendarHandler(calendarService service.CalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService: calendarService}
}

func (h *CalendarHandler) GetEvents(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	events, err := h.calendarService.GetEvents(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve calendar events", nil)
		return
	}

	response.Success(c, http.StatusOK, "Calendar events retrieved successfully", events)
}

func (h *CalendarHandler) GetEventByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid event ID", nil)
		return
	}

	event, err := h.calendarService.GetEventByID(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve calendar event", nil)
		return
	}

	response.Success(c, http.StatusOK, "Calendar event retrieved successfully", event)
}

func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	event, err := h.calendarService.CreateEvent(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create calendar event", nil)
		return
	}

	response.Success(c, http.StatusCreated, "Calendar event created successfully", event)
}

func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid event ID", nil)
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	event, err := h.calendarService.UpdateEvent(c.Request.Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update calendar event", nil)
		return
	}

	response.Success(c, http.StatusOK, "Calendar event updated successfully", event)
}

func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid event ID", nil)
		return
	}

	err = h.calendarService.DeleteEvent(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete calendar event", nil)
		return
	}

	response.Success(c, http.StatusOK, "Calendar event deleted successfully", nil)
}
