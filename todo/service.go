// package todo provides an API to store and manage Todos.
package todo

import (
	"context"
	"errors"
	"sort"
	"sync"
)

type Service interface {
	// Todo fetches a single Todo using the given 'id'.
	Todo(ctx context.Context, id int) (Todo, error)

	// Todos fetches all Todos.
	Todos(ctx context.Context) ([]Todo, error)

	// CreateTodo creates a new Todo with the given 'description'. All new Todo
	// items are initially marked as incomplete.
	CreateTodo(ctx context.Context, description string) error

	// UpdateTodo modifies an existing Todo.
	UpdateTodo(ctx context.Context, id int, desc string) error

	// SetCompletedStatus sets a Todo item's completed status.
	SetCompletedStatus(ctx context.Context, id int, completed bool) error
}

// Todo represents a Todo task.
type Todo struct {
	ID          int
	Description string
	Completed   bool
}

type inmemService struct {
	nextID int

	todos map[int]Todo
	mu    sync.Mutex
}

// NewInmemService returns an in memory implementation of Service.
func NewInmemService() Service {
	return &inmemService{
		nextID: 1,
		todos:  make(map[int]Todo),
	}
}

// Errors returned by Service.
var (
	ErrInvalidID    = errors.New("invalid id, must be non-negative")
	ErrTodoNotFound = errors.New("todo not found")
)

// Todo implements Service.
func (s *inmemService) Todo(ctx context.Context, id int) (Todo, error) {
	if id < 0 {
		return Todo{}, ErrInvalidID
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if todo, ok := s.todos[id]; ok {
		return todo, nil
	}

	return Todo{}, nil
}

// Todos implements Service.
func (s *inmemService) Todos(ctx context.Context) ([]Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var todos []Todo
	for _, v := range s.todos {
		todos = append(todos, v)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})
	return todos, nil
}

// CreateTodo implements Service.
func (s *inmemService) CreateTodo(ctx context.Context, desc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo := Todo{
		ID:          s.generateID(),
		Description: desc,
	}

	s.todos[todo.ID] = todo
	return nil
}

func (s *inmemService) generateID() int {
	defer func() { s.nextID++ }()
	return s.nextID
}

// UpdateTodo implements Service.
func (s *inmemService) UpdateTodo(ctx context.Context, id int, desc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return ErrTodoNotFound
	}

	todo.Description = desc
	s.todos[id] = todo

	return nil
}

// SetCompletedStatus implements Service.
func (s *inmemService) SetCompletedStatus(ctx context.Context, id int, completed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return ErrTodoNotFound
	}

	todo.Completed = completed

	s.todos[id] = todo

	return nil
}
