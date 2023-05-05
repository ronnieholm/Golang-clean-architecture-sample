package infrastructure

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ronnieholm/Golang-clean-architecture-sample/application"
	"github.com/ronnieholm/Golang-clean-architecture-sample/domain/storyAggregate"
)

type SystemClock struct {
}

// singletons (ensure thread-safe)
var clock = SystemClock{}

func (s SystemClock) Utc() time.Time {
	return time.Now()
}

type AppEnv struct {
	// scoped (ensure thread-safe)
	tx              *sql.Tx
	storyRepository StoryRepository
}

// TODO: create mock AppEnv for unit testing

func New() AppEnv {
	db, err := sql.Open("sqlite3", "../scrum.sqlite")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	return AppEnv{
		tx:              tx,
		storyRepository: NewStoryRepository(tx),
	}
}

func (e AppEnv) SystemClock() application.SystemClock {
	return clock
}

func (e AppEnv) StoryRepository() storyAggregate.StoryRepository {
	return e.storyRepository
}

func (e AppEnv) Commit() {
	if err := e.tx.Commit(); err != nil {
		panic(err)
	}
}

func (e AppEnv) Rollback() {
	if err := e.tx.Rollback(); err != nil {
		panic(err)
	}
}
