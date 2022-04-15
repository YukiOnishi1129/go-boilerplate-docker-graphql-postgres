package view

import (
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/timeutil"
	"strconv"
)

func NewTodoFromModel(entity *entity.Todo) *model.Todo {
	resUser := model.User{
		ID:        strconv.FormatUint(entity.R.User.ID, 10),
		Name:      entity.R.User.Name,
		Email:     entity.R.User.Email,
		CreatedAt: timeutil.TimeFormat(entity.R.User.CreatedAt),
		UpdatedAt: timeutil.TimeFormat(entity.R.User.UpdatedAt),
	}
	resTodo := model.Todo{
		ID:        strconv.FormatUint(entity.ID, 10),
		Title:     entity.Title,
		Comment:   entity.Comment,
		CreatedAt: timeutil.TimeFormat(entity.CreatedAt),
		UpdatedAt: timeutil.TimeFormat(entity.UpdatedAt),
	}

	if entity.R.User.ImageURL.Valid {
		imageURL := entity.R.User.ImageURL.String
		resUser.ImageURL = &imageURL
	}

	if entity.R.User.DeletedAt.Valid {
		userDeletedAt := timeutil.TimeFormat(entity.R.User.DeletedAt.Time)
		resUser.DeletedAt = &userDeletedAt
	}

	resTodo.User = &resUser

	if entity.DeletedAt.Valid {
		deletedAt := timeutil.TimeFormat(entity.DeletedAt.Time)
		resTodo.DeletedAt = &deletedAt
	}

	return &resTodo
}
