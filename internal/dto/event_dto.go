package dto

import "time"

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	Location    *string   `json:"location"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	Location    *string   `json:"location"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}

type EventResponse struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	GoogleEventID *string   `json:"google_event_id,omitempty"`
	Title         string    `json:"title"`
	Description   *string   `json:"description,omitempty"`
	Location      *string   `json:"location,omitempty"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
