package todo

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

var errBadRouting = errors.New("bad routing")

type idRequest struct {
	id int
}

type todoResponse struct {
	Todo Todo  `json:"todo"`
	Err  error `json:"error,omitempty"`
}

func (r todoResponse) error() error { return r.Err }

func makeTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(idRequest)
		if !ok {
			return nil, errBadRouting
		}

		todo, err := s.Todo(ctx, req.id)
		return todoResponse{Todo: todo, Err: err}, nil

	}
}

type todosResponse struct {
	Todos []Todo `json:"todos,omitempty"`
	Err   error  `json:"error,omitempty"`
}

func (r todosResponse) error() error { return r.Err }

func makeTodosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if request != nil {
			return nil, errBadRouting
		}

		todos, err := s.Todos(ctx)
		return todosResponse{Todos: todos, Err: err}, nil
	}
}

type createTodoRequest struct {
	Description string `json:"description"`
}

type errorResponse struct {
	Err error `json:"-"`
}

func (r errorResponse) error() error { return r.Err }

func makeCreateTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(createTodoRequest)
		if !ok {
			return nil, errBadRouting
		}

		err := s.CreateTodo(ctx, req.Description)
		return errorResponse{Err: err}, nil
	}
}

type updateTodoRequest struct {
	ID          int    `json:"-"`
	Description string `json:"description"`
}

func makeUpdateTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(updateTodoRequest)
		if !ok {
			return nil, errBadRouting
		}

		err := s.UpdateTodo(ctx, req.ID, req.Description)
		return errorResponse{Err: err}, nil

	}
}

type setCompletedStatusRequest struct {
	ID        int  `json:"-"`
	Completed bool `json:"completed"`
}

func makeSetCompletedStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(setCompletedStatusRequest)
		if !ok {
			return nil, errBadRouting
		}

		err := s.SetCompletedStatus(ctx, req.ID, req.Completed)
		return errorResponse{Err: err}, nil
	}
}
