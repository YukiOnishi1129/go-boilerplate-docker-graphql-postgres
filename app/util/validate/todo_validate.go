package validate

import (
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

func CreateTodoValidation(input model.CreateTodoInput) error {
	return validation.ValidateStruct(&input,
		validation.Field(
			&input.Title,
			validation.Required.Error("タイトルは必須入力です。"),
			validation.RuneLength(1, 10).Error("タイトルは 1～10 文字です"),
		),
		validation.Field(
			&input.Comment,
			validation.Required.Error("コメントは必須入力です。"),
		),
	)
}

func UpdateTodoValidation(input model.UpdateTodoInput) error {
	return validation.ValidateStruct(&input,
		validation.Field(
			&input.ID,
			validation.Required.Error("IDは必須です。"),
		),
		validation.Field(
			&input.Title,
			validation.Required.Error("タイトルは必須入力です。"),
			validation.RuneLength(1, 10).Error("タイトルは 1～10 文字です"),
		),
		validation.Field(
			&input.Comment,
			validation.Required.Error("コメントは必須入力です。"),
		),
	)
}
