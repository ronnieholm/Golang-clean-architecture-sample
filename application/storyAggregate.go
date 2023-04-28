package application

import (
	"errors"
	"fmt"
	"time"

	"example.com/m/domain/storyAggregate"
	"github.com/google/uuid"
)

type CreateStoryCommand struct {
	Id          uuid.UUID
	Title       string
	Description *string
}

func (c CreateStoryCommand) Run(env AppEnv) (*uuid.UUID, error) {
	id := storyAggregate.NewStoryId(c.Id)
	duplicate, err := env.StoryRepository().Exist(id)
	if err != nil {
		return nil, fmt.Errorf("cannot query for existing story: %v", err)
	}
	if *duplicate {
		return nil, NewDuplicateError(id.Value)
	}

	title, titleErr := storyAggregate.NewStoryTitle("title", c.Title)
	var description *storyAggregate.StoryDescription
	var descriptionErr error
	if c.Description != nil {
		description, descriptionErr = storyAggregate.NewStoryDescription("description", *c.Description)
	}

	validationErrors := errors.Join(titleErr, descriptionErr)
	if validationErrors != nil {
		return nil, fmt.Errorf("validation errors: %v", validationErrors)
	}

	story, err := storyAggregate.NewStory(id, *title, description, env.SystemClock().Utc())
	if err != nil {
		return nil, fmt.Errorf("validation errors creating story: %v", validationErrors)
	}

	err = env.StoryRepository().ApplyEvent(*story)
	if err != nil {
		return nil, fmt.Errorf("cannot apply events: %v", err)
	}

	return (*uuid.UUID)(&story.Id.Value), nil
}

type AddTaskToStoryCommand struct {
	TaskId      uuid.UUID
	StoryId     uuid.UUID
	Title       string
	Description *string
}

func (c AddTaskToStoryCommand) Run(env AppEnv) (*uuid.UUID, error) {
	storyId := storyAggregate.NewStoryId(c.StoryId)
	story, err := env.StoryRepository().GetById(storyId)
	if err != nil {
		return nil, fmt.Errorf("story not found: %v", err)
	}

	taskId := storyAggregate.NewTaskId(c.TaskId)
	title, titleErr := storyAggregate.NewTaskTitle("title", c.Title)
	var description *storyAggregate.TaskDescription
	var descriptionErr error
	if c.Description != nil {
		description, descriptionErr = storyAggregate.NewTaskDescription("description", *c.Description)
	}

	validationErrors := errors.Join(titleErr, descriptionErr)
	if validationErrors != nil {
		return nil, fmt.Errorf("validation errors: %v", validationErrors)
	}

	task, err := storyAggregate.NewTask(taskId, *title, description, env.SystemClock().Utc())
	if err != nil {
		return nil, fmt.Errorf("validation error create task: %v", err)
	}

	err = story.AddTask(*task)
	if err != nil {
		return nil, fmt.Errorf("cannot add task: %v", err)
	}

	err = env.StoryRepository().ApplyEvent(*story)
	if err != nil {
		return nil, fmt.Errorf("cannot apply events: %v", err)
	}

	return (*uuid.UUID)(&taskId.Value), nil
}

type TaskDto struct {
	Id          uuid.UUID
	Title       string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewTaskDto(t storyAggregate.Task) TaskDto {
	return TaskDto{
		Id:          t.Id.Value,
		Title:       t.Title.Value,
		Description: nil,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

type StoryDto struct {
	Id          uuid.UUID
	Title       string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Tasks       []TaskDto
}

func NewStoryDto(s storyAggregate.Story) StoryDto {
	tasks := make([]TaskDto, len(s.Tasks))
	for _, t := range s.Tasks {
		tasks = append(tasks, NewTaskDto(t))
	}

	return StoryDto{
		Id:          s.Id.Value,
		Title:       s.Title.Value,
		Description: nil,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Tasks:       tasks,
	}
}

type GetStoryByIdQuery struct {
	Id uuid.UUID
}

func (q GetStoryByIdQuery) Run(env AppEnv) (*StoryDto, error) {
	storyId := storyAggregate.NewStoryId(q.Id)
	story, err := env.StoryRepository().GetById(storyId)
	if err != nil {
		return nil, fmt.Errorf("story not found: %v", err)
	}
	dto := NewStoryDto(*story)
	return &dto, nil
}
