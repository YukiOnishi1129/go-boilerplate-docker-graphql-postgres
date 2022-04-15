package view

import (
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/timeutil"
	"strconv"
)

func NewUserFromModel(entity *entity.User) *model.User {
	resUser := model.User{
		ID:        strconv.FormatUint(entity.ID, 10),
		Name:      entity.Name,
		Email:     entity.Email,
		CreatedAt: timeutil.TimeFormat(entity.CreatedAt),
		UpdatedAt: timeutil.TimeFormat(entity.UpdatedAt),
	}

	if entity.ImageURL.Valid {
		imageURL := entity.ImageURL.String
		resUser.ImageURL = &imageURL
	}

	if entity.DeletedAt.Valid {
		deletedAt := timeutil.TimeFormat(entity.DeletedAt.Time)
		resUser.DeletedAt = &deletedAt
	}

	return &resUser
}
