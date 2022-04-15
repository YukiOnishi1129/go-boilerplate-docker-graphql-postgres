package awsutil

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/view"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type AwsUtil struct {
	session *session.Session
}

func LazyInitTodoService(session *session.Session) *AwsUtil {
	return &AwsUtil{
		session,
	}
}

//Init 初期設定 aws sessionの作成
func Init() (*session.Session, error) {
	// sessionの作成
	sess, sessionErr := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			}),
			Region:                        aws.String(os.Getenv("AWS_S3_REGION")),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		Profile: os.Getenv("AWS_CREDENTIAL_USER"),
	})
	if sessionErr != nil {
		return nil, view.NewInternalServerErrorFromModel(sessionErr.Error())
	}

	return sess, nil
}

//UploadImageFileToS3 画像ファイルアップロード処理
func (a *AwsUtil) UploadImageFileToS3(file *graphql.Upload) (string, error) {
	filePath := fmt.Sprintf("images/%s", file.Filename)
	// S3へアップロード
	uploader := s3manager.NewUploader(a.session)
	_, uploadErr := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:         aws.String(filePath),
		Body:        file.File,
		ContentType: aws.String(file.ContentType),
	})
	if uploadErr != nil {
		return "", uploadErr
	}
	// cloud front経由の画像urlを作成
	imagePath := fmt.Sprintf("%s/%s", os.Getenv("AWS_CLOUD_FRONT_URL"), filePath)

	return imagePath, nil
}
