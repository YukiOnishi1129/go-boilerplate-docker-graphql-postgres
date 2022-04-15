package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/auth"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/view"
)

func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.User, error) {
	return r.userService.SignIn(ctx, input)
}

func (r *mutationResolver) SignUp(ctx context.Context, input model.SignUpInput) (*model.User, error) {
	return r.userService.SignUp(ctx, input)
}

func (r *mutationResolver) SignOut(ctx context.Context) (string, error) {
	return r.userService.SignOut(ctx)
}

func (r *mutationResolver) UpdateUserName(ctx context.Context, name string) (*model.User, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.UpdateUserName(ctx, name, adminUser)
}

func (r *mutationResolver) UpdateUserEmail(ctx context.Context, email string) (*model.User, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.UpdateUserEmail(ctx, email, adminUser)
}

func (r *mutationResolver) UpdateUserPassword(ctx context.Context, input model.UpdatePasswordInput) (*model.User, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.UpdateUserPassword(ctx, input, adminUser)
}

func (r *mutationResolver) UploadUserFile(ctx context.Context, file *graphql.Upload) (*model.User, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.UploadUserFile(ctx, file, adminUser)
}

func (r *queryResolver) MyUserDetail(ctx context.Context) (*model.User, error) {
	adminUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.MyUserDetail(adminUser)
}

func (r *queryResolver) UserDetail(ctx context.Context, id string) (*model.User, error) {
	_, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, view.NewUnauthorizedErrorFromModel(err.Error())
	}
	return r.userService.UserDetail(ctx, id)
}
