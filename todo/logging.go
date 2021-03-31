package todo

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	s      Service
	logger log.Logger
}

func NewLoggingService(s Service, logger log.Logger) Service {
	return &loggingService{
		s:      s,
		logger: logger,
	}
}

func (s *loggingService) Todo(ctx context.Context, id int) (todo Todo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"event", "get_todo",
			"id", id,
			"dur", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.Todo(ctx, id)
}

func (s *loggingService) Todos(ctx context.Context) (todos []Todo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"event", "get_todos",
			"count", len(todos),
			"dur", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.Todos(ctx)
}

func (s *loggingService) CreateTodo(ctx context.Context, desc string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"event", "post_create_todo",
			"dur", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.CreateTodo(ctx, desc)
}

func (s *loggingService) UpdateTodo(ctx context.Context, id int, desc string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"event", "put_update_todo",
			"id", id,
			"dur", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.UpdateTodo(ctx, id, desc)
}

func (s *loggingService) SetCompletedStatus(ctx context.Context, id int, completed bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"event", "put_set_todo_completed",
			"id", id,
			"completed", completed,
			"dur", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.SetCompletedStatus(ctx, id, completed)
}
