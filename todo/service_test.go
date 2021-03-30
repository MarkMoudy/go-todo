package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTodo(t *testing.T) {
	testcases := []struct {
		name     string
		todos    []Todo
		id       int
		wantTodo Todo
		errMatch error
	}{
		{
			name: "successful-lookup",
			todos: []Todo{
				{
					Description: "do this task",
				},
			},
			id: 1,
			wantTodo: Todo{
				ID:          1,
				Description: "do this task",
			},
		},
		{
			name:     "invalid-id",
			id:       -1,
			errMatch: ErrInvalidID,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewInmemService()
			ctx := context.Background()

			for _, v := range tc.todos {
				if err := svc.CreateTodo(ctx, v.Description); err != nil {
					t.Fatalf("unexpected error creating initial set of Todo data %v", err)
				}
			}

			gotTodo, gotErr := svc.Todo(ctx, tc.id)

			if !errors.Is(gotErr, tc.errMatch) {
				t.Fatalf("unexpected error got %v, want %v", gotErr, tc.errMatch)
			}

			if got, want := gotTodo, tc.wantTodo; got != want {
				t.Errorf("unexpected Todo returned got %+v, want %+v", got, want)
			}

		})
	}
}

func TestTodos(t *testing.T) {
	testcases := []struct {
		name      string
		todos     []Todo
		wantTodos []Todo
	}{
		{
			name: "no-items",
		},
		{
			name: "single-item",
			todos: []Todo{
				{
					Description: "desc 1",
				},
			},
			wantTodos: []Todo{
				{
					ID:          1,
					Description: "desc 1",
				},
			},
		},
		{
			name: "multiple-items",
			todos: []Todo{
				{
					Description: "desc 1",
				},
				{
					Description: "desc 2",
				},
				{
					Description: "desc 3",
				},
			},
			wantTodos: []Todo{
				{
					ID:          1,
					Description: "desc 1",
				},
				{
					ID:          2,
					Description: "desc 2",
				},
				{
					ID:          3,
					Description: "desc 3",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewInmemService()
			ctx := context.Background()

			for _, v := range tc.todos {
				if err := svc.CreateTodo(ctx, v.Description); err != nil {
					t.Fatalf("unexpected error inserting todo data %v", err)
				}
			}

			gotTodos, err := svc.Todos(ctx)
			if err != nil {
				t.Fatalf("unexpected error fetching Todos %v", err)
			}

			if got, want := gotTodos, tc.wantTodos; !cmp.Equal(got, want) {
				t.Errorf("unexpected todos returned (got: -, want: +) %v", cmp.Diff(got, want))
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	testcases := []struct {
		name     string
		todo     Todo
		desc     string
		wantTodo Todo
		errMatch error
	}{
		{
			name: "successful-update",
			todo: Todo{
				Description: "foo desc",
			},
			desc: "bar desc",
			wantTodo: Todo{
				ID:          1,
				Description: "bar desc",
			},
		},
		{
			name:     "todo-not-found",
			errMatch: ErrTodoNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewInmemService()
			ctx := context.Background()

			if err := svc.CreateTodo(ctx, tc.todo.Description); err != nil {
				t.Fatalf("unexpected error creating initial Todo %v", err)
			}

			if err := svc.UpdateTodo(ctx, tc.wantTodo.ID, tc.desc); !errors.Is(tc.errMatch, err) {
				t.Fatalf("unexpected error got %v, want %v", err, tc.errMatch)
			}

			gotTodo, err := svc.Todo(ctx, tc.wantTodo.ID)
			if err != nil {
				t.Fatalf("unexpected error fetching updated Todo %v", err)
			}

			if got, want := gotTodo, tc.wantTodo; !cmp.Equal(got, want) {
				t.Errorf("unexpected Todo returned (got: -, want: +) %v", cmp.Diff(got, want))
			}
		})
	}
}
