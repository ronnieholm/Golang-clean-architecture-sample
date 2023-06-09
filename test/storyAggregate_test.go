package infrastrucre

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ronnieholm/Golang-clean-architecture-sample/application"
	"github.com/ronnieholm/Golang-clean-architecture-sample/infrastructure"
	"github.com/stretchr/testify/assert"
)

func BuildCreateStoryCommand() application.CreateStoryCommand {
	return application.CreateStoryCommand{
		Id:          uuid.New(),
		Title:       "Title",
		Description: nil,
	}
}

func BuildAddTaskToStoryCommand(storyId uuid.UUID) application.AddTaskToStoryCommand {
	return application.AddTaskToStoryCommand{
		TaskId:      uuid.New(),
		StoryId:     storyId,
		Title:       "Title",
		Description: nil,
	}
}

func TestCreateStoryCommand(t *testing.T) {
	env := infrastructure.New()
	cmd := BuildCreateStoryCommand()
	storyId, err := cmd.Run(env)
	assert.Nil(t, err)

	qry := application.GetStoryByIdQuery{
		Id: *storyId,
	}
	story, err := qry.Run(env)
	assert.Nil(t, err)
	assert.NotNil(t, story)
	env.Commit()
}

func TestAddTaskToStoryCommand(t *testing.T) {
	env := infrastructure.New()
	storyCmd := BuildCreateStoryCommand()
	storyId, err := storyCmd.Run(env)
	assert.Nil(t, err)

	taskCmd := BuildAddTaskToStoryCommand(*storyId)
	taskId, err := taskCmd.Run(env)
	_ = taskId
	_ = err

	qry := application.GetStoryByIdQuery{
		Id: *storyId,
	}
	story, err := qry.Run(env)
	assert.Nil(t, err)
	assert.NotNil(t, story)
	env.Commit()
}
