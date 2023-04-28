package storyAggregate

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Value objects

func TestNewStoryTitle(t *testing.T) {
	// TODO: test field name in error
	var tests = []struct {
		input   string
		success bool
	}{
		{"", false},
		{" ", false},
		{"Title", true},
		{strings.Repeat("x", 51), false},
	}

	for _, tt := range tests {
		title, err := NewStoryTitle("nil", tt.input)
		if !tt.success && err == nil {
			t.Errorf("validation error expected")
		}
		if tt.success {
			assert.Equal(t, tt.input, title.Value)
		}
	}
}

// Entities

func TestNewStory(t *testing.T) {
	// TODO: implement StoryBuilder (can be used in integration test too)
	id := NewStoryId(uuid.New())
	title, _ := NewStoryTitle("nil", "Title")
	description, _ := NewStoryDescription("nil", "Description")
	createdAt := time.Now()
	story, err := NewStory(id, *title, description, createdAt)

	assert.NotNil(t, story)
	//assert.NotNil(t, event)
	assert.Nil(t, err)
}

func TestAddTask(t *testing.T) {
	// TODO: implement TaskBuilder (can be used in integration test too)
	id := NewStoryId(uuid.New())
	title, _ := NewStoryTitle("nil", "Title")
	description, _ := NewStoryDescription("nil", "Description")
	createdAt := time.Now()
	story, _ := NewStory(id, *title, description, createdAt)

	id2 := NewTaskId(uuid.New())
	title2, _ := NewTaskTitle("nil", "Title")
	description2, _ := NewTaskDescription("nil", "Description")
	task, _ := NewTask(id2, *title2, description2, time.Now())

	err := story.AddTask(*task)
	assert.Len(t, story.Tasks, 1)
	//assert.NotNil(t, event)
	assert.Nil(t, err)

	err2 := story.AddTask(*task)
	//assert.Nil(t, event2)
	assert.NotNil(t, err2)
}
