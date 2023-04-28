package infrastrucre

import (
	"testing"

	"example.com/m/application"
	"example.com/m/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func BuildCreateStoryCommandBuilder() application.CreateStoryCommand {
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
	cmd := BuildCreateStoryCommandBuilder()
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
	storyCmd := BuildCreateStoryCommandBuilder()
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
