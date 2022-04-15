package graph

import (
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService *service.UserService
	todoService *service.TodoService
}

func NewResolver(
	userService *service.UserService,
	todoService *service.TodoService,
) *Resolver {
	return &Resolver{
		userService: userService,
		todoService: todoService,
	}
}
