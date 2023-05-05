package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ronnieholm/Golang-clean-architecture-sample/domain/storyAggregate"
)

type StoryRepository struct {
	tx *sql.Tx
}

func NewStoryRepository(tx *sql.Tx) StoryRepository {
	return StoryRepository{
		tx: tx,
	}
}

func (r StoryRepository) Exist(id storyAggregate.StoryId) (*bool, error) {
	// TODO: prepare statement outside of any function for reuse?
	stmt, err := r.tx.Prepare("select count(*) from stories where id = ?")
	if err != nil {
		panic(err)
	}

	var x int64
	err = stmt.QueryRow(id.Value.String()).Scan(&x)
	if err != nil {
		panic(err)
	}

	exist := x == 1
	return &exist, nil
}

func (r StoryRepository) GetById(id storyAggregate.StoryId) (*storyAggregate.Story, error) {
	stmt, err := r.tx.Prepare("select id, title, description, created_at, updated_at from stories where id = ?")
	if err != nil {
		panic(err)
	}

	// https://www.calhoun.io/querying-for-multiple-records-with-gos-sql-package/
	// TODO: Create join query and parse: https://stackoverflow.com/questions/54601529/efficiently-mapping-one-to-many-many-to-many-database-to-struct-in-golang (handle zero tasks, left join?)

	var id2 string
	var title string
	var description sql.NullString
	var createdAt string
	var updatedAt sql.NullString

	err = stmt.QueryRow(id.Value.String()).Scan(&id2, &title, &description, &createdAt, &updatedAt)
	if err != nil {
		panic(err)
	}

	var story = storyAggregate.Story{
		Id:    storyAggregate.StoryId{Value: uuid.New()},
		Title: storyAggregate.StoryTitle{Value: title},
		// TODO: remaining fields
	}

	return &story, nil
}

func (r StoryRepository) ApplyEvent(s storyAggregate.Story) error {
	for len(s.Events) > 0 {
		switch v := s.Events[0].(type) {
		case storyAggregate.StoryCreatedEvent:
			// TODO: use defer to close statements?
			stmt, err := r.tx.Prepare("insert into stories (id, title, description, created_at) values (?, ?, ?, ?)")
			if err != nil {
				panic(err)
			}
			_, err = stmt.Exec(
				v.StoryId.Value.String(),
				v.StoryTitle.Value,
				nil,
				v.OccurredAt.String())
			if err != nil {
				panic(err)
			}
		case storyAggregate.TaskAddedToStoryEvent:
			stmt, err := r.tx.Prepare("insert into tasks (id, story_id, title, description, created_at) values (?, ?, ?, ?, ?)")
			if err != nil {
				panic(err)
			}
			_, err = stmt.Exec(
				v.TaskId.Value.String(),
				v.StoryId.Value.String(),
				v.TaskTitle.Value,
				nil,
				v.OccurredAt.String())
			if err != nil {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("%T", v))
		}
		s.Events = s.Events[1:]
	}

	return nil
}
