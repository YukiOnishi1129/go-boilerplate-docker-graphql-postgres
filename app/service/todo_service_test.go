package service

import (
	"context"
	"database/sql"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
	"testing"
)

const TimeLayout = "2006-01-02 15:04:05"

func TestService_TodoList_OnSuccess(t *testing.T) {
	RunWithDB(t, "get TodoList", func(t *testing.T, db *sql.DB) {
		//　予測値
		wantUser := model.User{
			ID:        strconv.FormatUint(1, 10),
			Name:      "太郎",
			Email:     "taro@gmail.com",
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		want := [...]*model.Todo{
			{
				ID:        strconv.FormatUint(1, 10),
				Title:     "todo1",
				Comment:   "todo1のコメント",
				User:      &wantUser,
				CreatedAt: TimeLayout,
				UpdatedAt: TimeLayout,
			},
			{
				ID:        strconv.FormatUint(2, 10),
				Title:     "todo2",
				Comment:   "todo2のコメント",
				User:      &wantUser,
				CreatedAt: TimeLayout,
				UpdatedAt: TimeLayout,
			},
		}

		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		//	実行
		result, resErr := s.TodoList(context.Background(), targetUser)
		if resErr != nil {
			t.Errorf("get TodoList() error = %v", resErr)
		}

		// テスト結果の評価
		for i, res := range result {
			if diff := cmp.Diff(*res, *want[i], cmpopts.IgnoreFields(*res, "CreatedAt", "UpdatedAt", "DeletedAt", "User.CreatedAt", "User.UpdatedAt", "User.DeletedAt")); diff != "" {
				t.Errorf("%v", diff)
			}
		}
	})
}

func TestService_TodoDetail_OnSuccess(t *testing.T) {
	RunWithDB(t, "get TodoDetail", func(t *testing.T, db *sql.DB) {
		//　予測値
		wantUser := model.User{
			ID:        strconv.FormatUint(1, 10),
			Name:      "太郎",
			Email:     "taro@gmail.com",
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		//　予測値
		want := model.Todo{
			ID:        strconv.FormatUint(2, 10),
			Title:     "todo2",
			Comment:   "todo2のコメント",
			User:      &wantUser,
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		targetID := 2
		//	実行
		result, resErr := s.TodoDetail(context.Background(), strconv.Itoa(targetID), targetUser)
		if resErr != nil {
			t.Errorf("get TodoDetail() error = %v", resErr)
		}

		if diff := cmp.Diff(*result, want, cmpopts.IgnoreFields(*result, "CreatedAt", "UpdatedAt", "DeletedAt", "User.CreatedAt", "User.UpdatedAt", "User.DeletedAt")); diff != "" {
			t.Errorf("%v", diff)
		}
	})
}

func TestService_TodoDetail_OnFailure(t *testing.T) {
	RunWithDB(t, "get TodoDetail error", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		targetID := 4
		//	実行
		result, resErr := s.TodoDetail(context.Background(), strconv.Itoa(targetID), targetUser)

		if resErr == nil {
			t.Fatalf("存在しないtodoはエラーになるべきです. err: %v", resErr)
		}
		if result != nil {
			t.Errorf("nilであるべきです. got: %v", result)
		}
	})
}

func TestService_CreateTodo_OnSuccess(t *testing.T) {
	RunWithDB(t, "create todo success", func(t *testing.T, db *sql.DB) {
		//　予測値
		wantUser := model.User{
			ID:        strconv.FormatUint(1, 10),
			Name:      "太郎",
			Email:     "taro@gmail.com",
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		// 予測値
		want := model.Todo{
			ID:        strconv.FormatUint(4, 10),
			Title:     "todo4",
			Comment:   "todo4のコメント",
			User:      &wantUser,
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}

		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		args := model.CreateTodoInput{
			Title:   "todo4",
			Comment: "todo4のコメント",
		}
		//	実行
		result, resErr := s.CreateTodo(context.Background(), args, targetUser)
		if resErr != nil {
			t.Errorf("CreateTodo() error = %v", resErr)
		}
		if diff := cmp.Diff(*result, want, cmpopts.IgnoreFields(*result, "CreatedAt", "UpdatedAt", "DeletedAt", "User.CreatedAt", "User.UpdatedAt", "User.DeletedAt")); diff != "" {
			t.Errorf("%v", diff)
		}
	})
}

func TestService_CreateTodo_OnFailure(t *testing.T) {
	RunWithDB(t, "create todo bad request empty title", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		// 予測値
		s := &TodoService{
			db: db,
		}
		args := model.CreateTodoInput{
			Title:   "",
			Comment: "todo4のコメント",
		}
		//	実行
		result, resErr := s.CreateTodo(context.Background(), args, targetUser)

		if resErr == nil {
			t.Fatalf("titleのバリデーションエラーになるべきです. err: %v", resErr)
		}
		if result != nil {
			t.Errorf("nilであるべきです. got: %v", result)
		}
	})

	RunWithDB(t, "create todo bad request empty comment", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		// 予測値
		s := &TodoService{
			db: db,
		}
		args := model.CreateTodoInput{
			Title:   "todo4のタイトル",
			Comment: "",
		}
		//	実行
		result, resErr := s.CreateTodo(context.Background(), args, targetUser)

		if resErr == nil {
			t.Fatalf("commentのバリデーションエラーになるべきです. err: %v", resErr)
		}
		if result != nil {
			t.Errorf("nilであるべきです. got: %v", result)
		}
	})
}

func TestService_UpdateTodo_OnSuccess(t *testing.T) {
	RunWithDB(t, "update todo ", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 2)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		// 予測値
		wantUser := model.User{
			ID:        strconv.FormatUint(targetUser.ID, 10),
			Name:      targetUser.Name,
			Email:     targetUser.Email,
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		want := model.Todo{
			ID:        strconv.FormatUint(3, 10),
			Title:     "todo3title",
			Comment:   "todo3コメントupdate",
			User:      &wantUser,
			CreatedAt: TimeLayout,
			UpdatedAt: TimeLayout,
		}
		s := &TodoService{
			db: db,
		}
		args := model.UpdateTodoInput{
			ID:      "3",
			Title:   "todo3title",
			Comment: "todo3コメントupdate",
		}
		//	実行
		result, resErr := s.UpdateTodo(context.Background(), args, targetUser)
		if resErr != nil {
			t.Errorf("UpdateTodo() error = %v", resErr)
		}
		if diff := cmp.Diff(*result, want, cmpopts.IgnoreFields(*result, "CreatedAt", "UpdatedAt", "DeletedAt", "User.CreatedAt", "User.UpdatedAt", "User.DeletedAt")); diff != "" {
			t.Errorf("%v", diff)
		}
	})
}

func TestService_UpdateTodo_OnFailure(t *testing.T) {
	RunWithDB(t, "update todo bad request title empty", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 2)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		args := model.UpdateTodoInput{
			ID:      "3",
			Title:   "",
			Comment: "todo3コメントupdate",
		}
		//	実行
		result, resErr := s.UpdateTodo(context.Background(), args, targetUser)

		if resErr == nil {
			t.Fatalf("idのバリデーションエラーになるべきです. err: %v", resErr)
		}
		if result != nil {
			t.Errorf("nilであるべきです. got: %v", result)
		}
	})

	RunWithDB(t, "update todo not found", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 2)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		args := model.UpdateTodoInput{
			ID:      "4",
			Title:   "todo3title",
			Comment: "todo3タイトルupdate",
		}
		//	実行
		result, resErr := s.UpdateTodo(context.Background(), args, targetUser)

		if resErr == nil {
			t.Fatalf("該当データなしでエラーになるべきです. err: %v", resErr)
		}
		if result != nil {
			t.Errorf("nilであるべきです. got: %v", result)
		}
	})
}

func TestService_DeleteTodo_OnSuccess(t *testing.T) {
	RunWithDB(t, "delete todo ", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		// 予測値
		want := "1"
		s := &TodoService{
			db: db,
		}
		args := "1"
		//	実行
		result, resErr := s.DeleteTodo(context.Background(), args, targetUser)
		if resErr != nil {
			t.Errorf("UpdateTodo() error = %v", resErr)
		}
		if diff := cmp.Diff(result, want); diff != "" {
			t.Errorf("%v", diff)
		}
	})
}

func TestService_DeleteTodo_OnFailure(t *testing.T) {
	RunWithDB(t, "delete todo not empty id", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		args := ""
		//	実行
		result, resErr := s.DeleteTodo(context.Background(), args, targetUser)
		if resErr == nil {
			t.Fatalf("idのバリデーションエラーになるべきです. err: %v", resErr)
		}
		if result != "" {
			t.Errorf("空文字であるべきです. got: %v", result)
		}
	})

	RunWithDB(t, "delete todo bad not found", func(t *testing.T, db *sql.DB) {
		targetUser, userErr := entity.Users(qm.Where("id=?", 1)).One(context.Background(), db)
		if userErr != nil {
			t.Errorf("get TodoList() error = %v", userErr)
		}

		s := &TodoService{
			db: db,
		}
		args := "4"
		//	実行
		result, resErr := s.DeleteTodo(context.Background(), args, targetUser)
		if resErr == nil {
			t.Fatalf("該当データなしでエラーになるべきです. err: %v", resErr)
		}
		if result != "" {
			t.Errorf("空文字であるべきです. got: %v", result)
		}
	})
}
