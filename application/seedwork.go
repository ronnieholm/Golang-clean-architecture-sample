package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/ronnieholm/Golang-clean-architecture-sample/domain/storyAggregate"
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
