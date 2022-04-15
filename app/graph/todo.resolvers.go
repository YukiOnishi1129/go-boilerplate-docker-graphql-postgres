package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/auth"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/view"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.Todo, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.todoService.CreateTodo(ctx, input, adminUser)
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.UpdateTodoInput) (*model.Todo, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.todoService.UpdateTodo(ctx, input, adminUser)
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (string, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return "", view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.todoService.DeleteTodo(ctx, id, adminUser)
}

func (r *queryResolver) TodoList(ctx context.Context) ([]*model.Todo, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.todoService.TodoList(ctx, adminUser)
}

func (r *queryResolver) TodoDetail(ctx context.Context, id string) (*model.Todo, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.todoService.TodoDetail(ctx, id, adminUser)
}
