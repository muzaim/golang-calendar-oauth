package service

import (
	"context"
	"errors"
	"log"
	"time"

	"golang-test/internal/dto"
	"golang-test/internal/model"
	"golang-test/internal/repository"
	"golang-test/pkg/oauth"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type CalendarService interface {
	GetEvents(ctx context.Context, userID int64) ([]dto.EventResponse, error)
	GetEventByID(ctx context.Context, userID int64, id int64) (*dto.EventResponse, error)
	CreateEvent(ctx context.Context, userID int64, req dto.CreateEventRequest) (*dto.EventResponse, error)
	UpdateEvent(ctx context.Context, userID int64, id int64, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	DeleteEvent(ctx context.Context, userID int64, id int64) error
}

type calendarService struct {
	eventRepo   repository.EventRepository
	userRepo    repository.UserRepository
	googleOAuth *oauth.GoogleOAuth
}

func NewCalendarService(
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
	googleOAuth *oauth.GoogleOAuth,
) CalendarService {
	return &calendarService{
		eventRepo:   eventRepo,
		userRepo:    userRepo,
		googleOAuth: googleOAuth,
	}
}

func (s *calendarService) GetEvents(ctx context.Context, userID int64) ([]dto.EventResponse, error) {
	events, err := s.eventRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, s.toEventResponse(&event))
	}

	return responses, nil
}

func (s *calendarService) GetEventByID(ctx context.Context, userID int64, id int64) (*dto.EventResponse, error) {
	event, err := s.eventRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}

	resp := s.toEventResponse(event)
	return &resp, nil
}

func (s *calendarService) CreateEvent(ctx context.Context, userID int64, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	event := &model.Event{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	// Try to sync with Google Calendar API if user has connected Google account
	user, err := s.userRepo.FindByID(userID)
	if err == nil && user != nil && user.GoogleAccessToken != nil {
		srv, err := s.getGoogleCalendarClient(ctx, user)
		if err == nil && srv != nil {
			gEvent := &calendar.Event{
				Summary:     req.Title,
				Location:    getDereferencedString(req.Location),
				Description: getDereferencedString(req.Description),
				Start: &calendar.EventDateTime{
					DateTime: req.StartTime.Format(time.RFC3339),
				},
				End: &calendar.EventDateTime{
					DateTime: req.EndTime.Format(time.RFC3339),
				},
			}
			createdGEvent, err := srv.Events.Insert("primary", gEvent).Do()
			if err == nil && createdGEvent != nil {
				event.GoogleEventID = &createdGEvent.Id
			} else if err != nil {
				log.Printf("Google Calendar Sync Error on Create: %v", err)
			}
		}
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	resp := s.toEventResponse(event)
	return &resp, nil
}

func (s *calendarService) UpdateEvent(ctx context.Context, userID int64, id int64, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.eventRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}

	event.Title = req.Title
	event.Description = req.Description
	event.Location = req.Location
	event.StartTime = req.StartTime
	event.EndTime = req.EndTime

	// Try to update on Google Calendar if GoogleEventID exists
	if event.GoogleEventID != nil {
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user != nil && user.GoogleAccessToken != nil {
			srv, err := s.getGoogleCalendarClient(ctx, user)
			if err == nil && srv != nil {
				gEvent := &calendar.Event{
					Summary:     req.Title,
					Location:    getDereferencedString(req.Location),
					Description: getDereferencedString(req.Description),
					Start: &calendar.EventDateTime{
						DateTime: req.StartTime.Format(time.RFC3339),
					},
					End: &calendar.EventDateTime{
						DateTime: req.EndTime.Format(time.RFC3339),
					},
				}
				_, err = srv.Events.Update("primary", *event.GoogleEventID, gEvent).Do()
				if err != nil {
					log.Printf("Google Calendar Sync Error on Update: %v", err)
				}
			}
		}
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	resp := s.toEventResponse(event)
	return &resp, nil
}

func (s *calendarService) DeleteEvent(ctx context.Context, userID int64, id int64) error {
	event, err := s.eventRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotFound
	}

	// Try to delete from Google Calendar if synced
	if event.GoogleEventID != nil {
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user != nil && user.GoogleAccessToken != nil {
			srv, err := s.getGoogleCalendarClient(ctx, user)
			if err == nil && srv != nil {
				err := srv.Events.Delete("primary", *event.GoogleEventID).Do()
				if err != nil {
					log.Printf("Google Calendar Sync Error on Delete: %v", err)
				}
			}
		}
	}

	return s.eventRepo.Delete(id, userID)
}

func (s *calendarService) getGoogleCalendarClient(ctx context.Context, user *model.User) (*calendar.Service, error) {
	if user.GoogleAccessToken == nil {
		return nil, errors.New("no google access token")
	}

	var refreshToken string
	if user.GoogleRefreshToken != nil {
		refreshToken = *user.GoogleRefreshToken
	}

	var expiry time.Time
	if user.GoogleTokenExpiresAt != nil {
		expiry = *user.GoogleTokenExpiresAt
	}

	tok := &oauth2.Token{
		AccessToken:  *user.GoogleAccessToken,
		RefreshToken: refreshToken,
		Expiry:       expiry,
		TokenType:    "Bearer",
	}

	return s.googleOAuth.GetCalendarService(ctx, tok)
}

func (s *calendarService) toEventResponse(event *model.Event) dto.EventResponse {
	return dto.EventResponse{
		ID:            event.ID,
		UserID:        event.UserID,
		GoogleEventID: event.GoogleEventID,
		Title:         event.Title,
		Description:   event.Description,
		Location:      event.Location,
		StartTime:     event.StartTime,
		EndTime:       event.EndTime,
		CreatedAt:     event.CreatedAt,
		UpdatedAt:     event.UpdatedAt,
	}
}

func getDereferencedString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
