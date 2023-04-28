package storyAggregate

import (
	"fmt"
	"strings"
	"time"

	"example.com/m/domain"
	"github.com/google/uuid"
)

// Value objects

type StoryId struct {
	Value uuid.UUID
}

func NewStoryId(value uuid.UUID) StoryId {
	return StoryId{
		Value: value,
	}
}

type StoryTitle struct {
	Value string
}

func NewStoryTitle(field string, value string) (*StoryTitle, error) {
	if strings.TrimSpace(value) == "" {
		return nil, domain.NewValidationError(field, fmt.Sprintf("empty or whitespace: '%s'", value))
	}

	const maxLength = 50
	length := len(value)
	if length > maxLength {
		return nil, domain.NewValidationError(field, fmt.Sprintf("contains more than {%d} characters: '%s' is has length '%d'", maxLength, value, length))
	}

	return &StoryTitle{
		Value: value,
	}, nil
}

type StoryDescription struct {
	Value string
}

func NewStoryDescription(field string, value string) (*StoryDescription, error) {
	if strings.TrimSpace(value) == "" {
		return nil, domain.NewValidationError(field, fmt.Sprintf("empty or whitespace: '%s'", value))
	}

	const maxLength = 100
	length := len(value)
	if length > maxLength {
		return nil, domain.NewValidationError(field, fmt.Sprintf("contains more than {%d} characters: '%s' is has length '%d'", maxLength, value, length))
	}

	return &StoryDescription{
		Value: value,
	}, nil
}

type TaskId struct {
	Value uuid.UUID
}

func NewTaskId(value uuid.UUID) TaskId {
	return TaskId{
		Value: value,
	}
}

type TaskTitle struct {
	Value string
}

func NewTaskTitle(field string, value string) (*TaskTitle, error) {
	if strings.TrimSpace(value) == "" {
		return nil, domain.NewValidationError(field, fmt.Sprintf("empty or whitespace: '%s'", value))
	}

	const maxLength = 50
	length := len(value)
	if length > maxLength {
		return nil, domain.NewValidationError(field, fmt.Sprintf("contains more than {%d} characters: '%s' is has length '%d'", maxLength, value, length))
	}

	return &TaskTitle{
		Value: value,
	}, nil
}

type TaskDescription struct {
	Value string
}

func NewTaskDescription(field string, value string) (*TaskDescription, error) {
	if strings.TrimSpace(value) == "" {
		return nil, domain.NewValidationError(field, fmt.Sprintf("empty or whitespace: '%s'", value))
	}

	const maxLength = 100
	length := len(value)
	if length > maxLength {
		return nil, domain.NewValidationError(field, fmt.Sprintf("contains more than {%d} characters: '%s' is has length '%d'", maxLength, value, length))
	}

	return &TaskDescription{
		Value: value,
	}, nil
}

// Domain events

type DomainEvent struct {
	OccurredAt time.Time
}

type StoryCreatedEvent struct {
	DomainEvent
	StoryId          StoryId
	StoryTitle       StoryTitle
	StoryDescription *StoryDescription
}

type TaskAddedToStoryEvent struct {
	DomainEvent
	TaskId          TaskId
	StoryId         StoryId
	TaskTitle       TaskTitle
	TaskDescription *TaskDescription
}

// Entitities

type Task struct {
	domain.Entity
	Id          TaskId
	Title       TaskTitle
	Description *TaskDescription
}

func NewTask(id TaskId, title TaskTitle, description *TaskDescription, createdAt time.Time) (*Task, error) {
	return &Task{
		Entity: domain.Entity{
			CreatedAt: createdAt,
			UpdatedAt: nil,
		},
		Id:          id,
		Title:       title,
		Description: description,
	}, nil
}

type Story struct {
	domain.AggregateRoot
	Id          StoryId
	Title       StoryTitle
	Description *StoryDescription
	Tasks       []Task
}

func NewStory(id StoryId, title StoryTitle, description *StoryDescription, createdAt time.Time) (*Story, error) {
	story := &Story{
		AggregateRoot: domain.AggregateRoot{
			Entity: domain.Entity{
				CreatedAt: createdAt,
				UpdatedAt: nil,
			},
		},
		Id:          id,
		Title:       title,
		Description: description,
		Tasks:       []Task{},
	}
	event := StoryCreatedEvent{
		DomainEvent: DomainEvent{
			OccurredAt: createdAt,
		},
		StoryId:          id,
		StoryTitle:       title,
		StoryDescription: description,
	}
	story.Events = append(story.Events, event)
	return story, nil
}

func (s *Story) AddTask(task Task) error {
	for _, t := range s.Tasks {
		s1 := t.Id.Value.String()
		s2 := task.Id.Value.String()
		if s1 == s2 {
			return domain.NewValidationError("task", fmt.Sprintf("duplicate task with Id {%s}", task.Id.Value.String()))
		}
	}

	s.Tasks = append(s.Tasks, task)
	event := TaskAddedToStoryEvent{
		DomainEvent: DomainEvent{
			OccurredAt: task.CreatedAt,
		},
		TaskId:          task.Id,
		TaskTitle:       task.Title,
		TaskDescription: task.Description,
	}
	s.Events = append(s.Events, event)
	return nil
}

// Repositories

type StoryRepository interface {
	// TODO: add context argument
	Exist(id StoryId) (*bool, error)
	GetById(id StoryId) (*Story, error)
	ApplyEvent(s Story) error
}
