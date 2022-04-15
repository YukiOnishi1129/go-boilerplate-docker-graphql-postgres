package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/entity"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/pkg/errors"
)

type Admin struct {
	User *entity.User
}

type contextKey struct {
	uuid string
}

var (
	adminKey      = &contextKey{"admin"}
	httpWriterKey = &contextKey{"httpWriter"}
	authCookieKey = "auth-cookie"
)

// MiddleWare cookie認証
func MiddleWare(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 認証情報のcookieを取得
			c, authErr := r.Cookie(authCookieKey)
			// 書き込み用httpを設定 (これに認証用のcookieを設定する)
			ctx := context.WithValue(r.Context(), httpWriterKey, w)
			r = r.WithContext(ctx)
			if authErr != nil || c == nil {
				next.ServeHTTP(w, r)
				fmt.Printf("auth.middleware: %v\n", authErr)
				return
			}

			// cookieよりuserIdを取得
			userID, err := getUserIDFromJwt(c)
			if err != nil {
				next.ServeHTTP(w, r)
				fmt.Printf("auth.middleware get userID from cookie: %v\n", err)
				return
			}

			// userIdよりユーザー情報を取得
			user, userErr := entity.Users(qm.Where("id=$1", userID)).One(ctx, db)
			if userErr != nil {
				next.ServeHTTP(w, r)
				fmt.Printf("auth.middleware get user: %v\n", userErr)
				return
			}

			// contextにユーザー情報をセット
			ctx = context.WithValue(r.Context(), adminKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext contextに保持しているユーザー情報を取得する
func GetUserIDFromContext(ctx context.Context) (*entity.User, error) {
	adminUser, adminOk := ctx.Value(adminKey).(*entity.User)
	if !adminOk {
		return nil, errors.New("認証情報がありません。")
	}
	return adminUser, nil
}

// RemoveAuthCookie 期限なしのcookieを設定 (ログアウト用)
func RemoveAuthCookie(ctx context.Context) {
	writer, _ := ctx.Value(httpWriterKey).(http.ResponseWriter)
	http.SetCookie(writer, &http.Cookie{
		HttpOnly: true,
		MaxAge:   0,
		Secure:   true,
		Name:     authCookieKey,
	})
}

// SetAuthCookie 認証用cookieを設定
func SetAuthCookie(ctx context.Context, user *entity.User) {
	// jwt tokenを作成
	sessionToken, err := createJwtToken(user)
	if err != nil {
		fmt.Println("Error: create jwt error", err)
		return
	}

	writer, _ := ctx.Value(httpWriterKey).(http.ResponseWriter)

	// 有効期限は2週間
	week := 60 * 60 * 24 * 7

	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   week * 2,
		Secure:   true,
		Name:     authCookieKey,
		Value:    sessionToken,
	}
	http.SetCookie(writer, &cookie)
}

func createJwtToken(user *entity.User) (string, error) {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	//claims["sub"] = strconv.Itoa(int(user.ID)) + user.Email + user.Name
	claims["id"] = user.ID
	//claims["name"] = user.Name
	// latを取り除かないとミドルウェアで「Token used before issued」エラーになる
	// https://github.com/dgrijalva/jwt-go/issues/314#issuecomment-812775567
	// claims["iat"] = time.Now() // jwtの発行時間
	// 経過時間
	// 経過時間を過ぎたjetは処理しないようになる
	// ここでは2週間の経過時間をリミットにしている
	claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Unix()

	// 電子署名
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getUserIDFromJwt cookieからjwtトークンを取得し、そこからuserIDを取得する
func getUserIDFromJwt(c *http.Cookie) (int, error) {
	clientToken := c.Value
	if clientToken == "" {
		return 0, errors.New("not token")
	}

	secretKey := os.Getenv("JWT_KEY")

	token, err := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("トークンをjwtにparseできません。")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, claimOk := token.Claims.(jwt.MapClaims)
	if !claimOk || !token.Valid {
		return 0, errors.New("id type not match")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("id type not match")
	}

	return int(userID), nil
}
