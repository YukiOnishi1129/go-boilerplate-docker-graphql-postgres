package service

import (
	"context"
	"database/sql"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	validate "github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/validate"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/view"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TodoService struct {
	db *sql.DB
}

func LazyInitTodoService(db *sql.DB) *TodoService {
	return &TodoService{
		db: db,
	}
}

func (s *TodoService) TodoList(ctx context.Context, adminUser *entity.User) ([]*model.Todo, error) {
	todoList, todoErr := entity.Todos(qm.Where("user_id=?", adminUser.ID), qm.Load("User")).All(ctx, s.db)
	if todoErr != nil {
		return nil, view.NewDBErrorFromModel(todoErr)
	}
	resTodoList := make([]*model.Todo, len(todoList))
	for i, todo := range todoList {
		resTodoList[i] = view.NewTodoFromModel(todo)
	}
	return resTodoList, nil
}

func (s *TodoService) TodoDetail(ctx context.Context, id string, adminUser *entity.User) (*model.Todo, error) {
	// バリデーション
	if id == "" {
		return nil, view.NewBadRequestErrorFromModel("IDは必須です。")
	}
	todo, todoErr := entity.Todos(qm.Where("id=?", id), qm.Where("user_id=?", adminUser.ID), qm.Load("User")).One(ctx, s.db)
	if todoErr != nil {
		return nil, view.NewDBErrorFromModel(todoErr)
	}
	return view.NewTodoFromModel(todo), nil
}

func (s *TodoService) CreateTodo(ctx context.Context, input model.CreateTodoInput, adminUser *entity.User) (*model.Todo, error) {
	var err error
	// バリデーション
	if err = validate.CreateTodoValidation(input); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}

	// 新規登録処理
	newTodo := &entity.Todo{
		Title:   input.Title,
		Comment: input.Comment,
		UserID:  adminUser.ID,
	}
	if err = newTodo.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}

	todo, todoErr := entity.Todos(qm.Where("id=?", newTodo.ID), qm.Where("user_id=?", adminUser.ID), qm.Load("User")).One(ctx, s.db)
	if todoErr != nil {
		return nil, view.NewDBErrorFromModel(todoErr)
	}

	return view.NewTodoFromModel(todo), nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, input model.UpdateTodoInput, adminUser *entity.User) (*model.Todo, error) {
	var err error
	// バリデーション
	if err = validate.UpdateTodoValidation(input); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}
	todo, todoErr := entity.Todos(qm.Where("id=?", input.ID), qm.Where("user_id=?", adminUser.ID), qm.Load("User")).One(ctx, s.db)
	if todoErr != nil {
		return nil, view.NewDBErrorFromModel(todoErr)
	}

	// 更新処理
	todo.Title = input.Title
	todo.Comment = input.Comment
	_, err = todo.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}
	return view.NewTodoFromModel(todo), nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string, adminUser *entity.User) (string, error) {
	var err error
	// バリデーション
	if id == "" {
		return "", view.NewBadRequestErrorFromModel("IDは必須です。")
	}
	todo, todoErr := entity.Todos(qm.Where("id=?", id), qm.Where("user_id=?", adminUser.ID)).One(ctx, s.db)
	if todoErr != nil {
		return "", view.NewDBErrorFromModel(todoErr)
	}

	// 削除処置
	if _, err = todo.Delete(ctx, s.db); err != nil {
		return "", view.NewInternalServerErrorFromModel(todoErr.Error())
	}
	return id, nil
}
