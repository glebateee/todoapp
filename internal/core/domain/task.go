package domain

import (
	"fmt"
	"time"

	core_errors "github.com/glebateee/todoapp/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Time{},
		nil,
		authorUserID,
	)
	// return Task{
	// 	ID:           UninitializedID,
	// 	Version:      UninitializedVersion,
	// 	Title:        title,
	// 	Description:  description,
	// 	AuthorUserID: authorUserID,
	// }
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid title length: %d: %w", titleLen, core_errors.ErrInvalidArgument)
	}
	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid description length: %d: %w", descriptionLen, core_errors.ErrInvalidArgument)

		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("task completed but completion time not set: %w", core_errors.ErrInvalidArgument)

		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("completion can't be before creation time: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("task isn't completed but completion time set: %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}
