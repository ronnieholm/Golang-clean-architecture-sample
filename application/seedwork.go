package application

import (
	"time"

	"example.com/m/domain/storyAggregate"
	"github.com/google/uuid"
)

type DuplicateError struct {
	Id uuid.UUID
}

func NewDuplicateError(id uuid.UUID) DuplicateError {
	return DuplicateError{
		Id: id,
	}
}

func (e DuplicateError) Error() string {
	return "duplicate error string"
}

type SystemClock interface {
	Utc() time.Time
}

type AppEnv interface {
	SystemClock() SystemClock
	StoryRepository() storyAggregate.StoryRepository
	Commit()
	Rollback()
}

type NotificationHub struct{}
