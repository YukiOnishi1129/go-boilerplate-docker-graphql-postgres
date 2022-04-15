package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/graph/model"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/auth"
	awsutil "github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/aws"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/validate"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/view"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db      *sql.DB
	awsUtil *awsutil.AwsUtil
}

func LazyInitUserService(db *sql.DB, awsUtil *awsutil.AwsUtil) *UserService {
	return &UserService{
		db:      db,
		awsUtil: awsUtil,
	}
}

// SignIn ログイン
func (s *UserService) SignIn(ctx context.Context, input model.SignInInput) (*model.User, error) {
	var err error
	// バリデーション
	if err = validate.SignInValidation(input); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}

	// ユーザー認証
	user, err := entity.Users(qm.Where("email=?", input.Email)).One(ctx, s.db)
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}
	// パスワード照合
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, view.NewUnauthorizedErrorFromModel("パスワードが違います。")
	}

	// cookieに保持
	auth.SetAuthCookie(ctx, user)
	return view.NewUserFromModel(user), nil
}

// SignUp 会員登録
func (s *UserService) SignUp(ctx context.Context, input model.SignUpInput) (*model.User, error) {
	var err error
	// バリデーション
	if err = validate.SignUpValidation(input); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}
	// パスワード判定処理
	if input.Password != input.PasswordConfirm {
		return nil, view.NewBadRequestErrorFromModel("パスワードが一致しません。")
	}
	// メールアドレス判定処理
	sameEmailUser, sameEmailUserErr := entity.Users(qm.Where("email=?", input.Email)).One(ctx, s.db)
	if sameEmailUserErr != nil && sameEmailUserErr.Error() != "sql: no rows in result set" {
		return nil, view.NewBadRequestErrorFromModel(sameEmailUserErr.Error())
	}
	if sameEmailUser != nil {
		return nil, view.NewBadRequestErrorFromModel(fmt.Sprintf("メールアドレス「%s」は使用されています。", input.Email))
	}

	// パスワードハッシュ化
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// DB登録
	newUser := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashPassword),
	}
	if err = newUser.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}

	// cookieに保持
	auth.SetAuthCookie(ctx, newUser)

	return view.NewUserFromModel(newUser), nil
}

// SignOut ログアウト
func (s *UserService) SignOut(ctx context.Context) (string, error) {
	auth.RemoveAuthCookie(ctx)
	return "ログアウトしました。", nil
}

// UpdateUserName ユーザー名変更
func (s *UserService) UpdateUserName(ctx context.Context, name string, adminUser *entity.User) (*model.User, error) {
	var err error
	// バリデーション
	if err = validate.UpdateUserNameValidation(validate.UpdateUserNameInput{Name: name}); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}
	// 更新処理
	adminUser.Name = name
	_, err = adminUser.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}
	return view.NewUserFromModel(adminUser), nil
}

// UpdateUserEmail ユーザーメールアドレス変更処理
func (s *UserService) UpdateUserEmail(ctx context.Context, email string, adminUser *entity.User) (*model.User, error) {
	var err error
	// バリデーション
	if err = validate.UpdateUserEmailValidation(validate.UpdateUserEmailInput{Email: email}); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}
	// メールアドレス判定処理
	sameEmailUser, sameEmailUserErr := entity.Users(qm.Where("email=?", email)).One(ctx, s.db)
	if sameEmailUserErr != nil && sameEmailUserErr.Error() != "sql: no rows in result set" {
		return nil, view.NewBadRequestErrorFromModel(sameEmailUserErr.Error())
	}
	if sameEmailUser != nil {
		return nil, view.NewBadRequestErrorFromModel(fmt.Sprintf("メールアドレス「%s」は使用されています。", email))
	}
	// 更新処理
	adminUser.Email = email
	_, err = adminUser.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}
	return view.NewUserFromModel(adminUser), nil
}

// UpdateUserPassword ユーザーパスワード変更
func (s *UserService) UpdateUserPassword(ctx context.Context, input model.UpdatePasswordInput, adminUser *entity.User) (*model.User, error) {
	var err error
	// バリデーション
	if err = validate.UpdateUserPasswordValidation(input); err != nil {
		return nil, view.NewBadRequestErrorFromModel(err.Error())
	}
	// 現在のパスワード照合処理
	targetUser, targetUserErr := entity.Users(qm.Where("id=?", adminUser.ID)).One(ctx, s.db)
	if targetUserErr != nil {
		return nil, view.NewBadRequestErrorFromModel(targetUserErr.Error())
	}
	if err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(input.OldPassword)); err != nil {
		return nil, view.NewUnauthorizedErrorFromModel("現在のパスワードが違います。")
	}
	// パスワード判定処理
	if input.NewPassword != input.NewPasswordConfirm {
		return nil, view.NewBadRequestErrorFromModel("新しいパスワードと確認用パスワードが一致しません。")
	}
	// 更新処理
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	adminUser.Password = string(hashPassword)
	_, err = adminUser.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}

	return view.NewUserFromModel(adminUser), nil
}

// UploadUserFile ファイルアップロード
func (s *UserService) UploadUserFile(ctx context.Context, file *graphql.Upload, adminUser *entity.User) (*model.User, error) {
	// 画像ファイルアップロード処理
	imagePath, uploadErr := s.awsUtil.UploadImageFileToS3(file)
	if uploadErr != nil {
		return nil, view.NewInternalServerErrorFromModel(uploadErr.Error())
	}

	// ユーザー情報を更新
	adminUser.ImageURL = null.StringFromPtr(&imagePath)
	_, err := adminUser.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, view.NewDBErrorFromModel(err)
	}

	return view.NewUserFromModel(adminUser), nil
}

// MyUserDetail ログインユーザーの詳細情報
func (s *UserService) MyUserDetail(adminUser *entity.User) (*model.User, error) {
	return view.NewUserFromModel(adminUser), nil
}

// UserDetail ユーザー詳細情報
func (s *UserService) UserDetail(ctx context.Context, id string) (*model.User, error) {
	// バリデーション
	if id == "" {
		return nil, view.NewBadRequestErrorFromModel("IDは必須です。")
	}
	user, userErr := entity.Users(qm.Where("id=?", id)).One(ctx, s.db)
	if userErr != nil {
		return nil, view.NewDBErrorFromModel(userErr)
	}
	return view.NewUserFromModel(user), nil
}
