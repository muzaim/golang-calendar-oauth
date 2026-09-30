package repository

import (
	"errors"

	"golang-test/internal/model"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	FindByIDAndUserID(id int64, userID int64) (*model.Event, error)
	FindAllByUserID(userID int64) ([]model.Event, error)
	Update(event *model.Event) error
	Delete(id int64, userID int64) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByIDAndUserID(id int64, userID int64) (*model.Event, error) {
	var event model.Event
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindAllByUserID(userID int64) ([]model.Event, error) {
	var events []model.Event
	err := r.db.Where("user_id = ?", userID).Order("start_time ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id int64, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Event{}).Error
}
