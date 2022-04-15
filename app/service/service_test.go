package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/testdata"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var (
	resourceTest *dockertest.Resource
	poolTest     *dockertest.Pool
	connTest     *sql.DB
)

var tableSQLFileName = [...]string{"users", "todos"}
var unEnabledFkKeySQLFileName = "un_enabled"
var enabledFkKeySQLFileName = "enabled"

func TestMain(m *testing.M) {
	beforeAll()
	defer afterAll()
	m.Run()
}

func RunWithDB(t *testing.T, name string, f func(t *testing.T, db *sql.DB)) {
	beforeEach()
	// テスト実行
	t.Run(name, func(t *testing.T) {
		f(t, connTest)
	})
}

func beforeAll() {
	fmt.Println("beforeAll")
	var err error
	// コンテナ起動
	createContainer()
	// テーブル作成
	_, fileName, _, _ := runtime.Caller(0)
	err = connectDB()
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	for _, sqlFileName := range tableSQLFileName {
		if err = execSQLScript(fmt.Sprintf("%s/../testdata/sql/create/%s.sql", filepath.Dir(fileName), sqlFileName)); err != nil {
			log.Fatalf("%s, %v", fileName, err)
		}
	}
}

func beforeEach() {
	var err error
	_, fileName, _, _ := runtime.Caller(0)
	// 外部キーを無効化
	if err = execSQLScript(fmt.Sprintf("%s/../testdata/sql/fkkey/%s.sql", filepath.Dir(fileName), unEnabledFkKeySQLFileName)); err != nil {
		log.Fatalf("%s, %v", fileName, err)
	}
	// データ削除
	for _, sqlFileName := range tableSQLFileName {
		if err = execSQLScript(fmt.Sprintf("%s/../testdata/sql/truncate/%s.sql", filepath.Dir(fileName), sqlFileName)); err != nil {
			log.Fatalf("%s, %v", fileName, err)
		}
	}
	// 外部キーを有効化
	if err = execSQLScript(fmt.Sprintf("%s/../testdata/sql/fkkey/%s.sql", filepath.Dir(fileName), enabledFkKeySQLFileName)); err != nil {
		log.Fatalf("%s, %v", fileName, err)
	}
	// テストデータ作成
	if err = createTestData(); err != nil {
		return
	}
}

func afterAll() {
	fmt.Println("afterAll")
	// コンテナ停止
	closeContainer()
}

func createContainer() {
	var err error
	// 絶対パスを取得
	_, fileName, _, _ := runtime.Caller(0)

	// Dockerとの接続
	poolTest, err = dockertest.NewPool("")
	poolTest.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Dockerコンテナ起動時の細かいオプションを指定する
	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
		Mounts: []string{
			fmt.Sprintf("%s/../../mysql/mysql.cnf:/etc/mysql/conf.d/mysql.cnf", filepath.Dir(fileName)),
			fmt.Sprintf("%s/../../mysql/db:/docker-entrypoint-initdb.d", filepath.Dir(fileName)), // コンテナ起動時に実行するSQL
		},
		Cmd: []string{
			"mysqld",
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
		},
	}

	// コンテナを起動
	resourceTest, err = poolTest.RunWithOptions(runOptions)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
}

func closeContainer() {
	var err error
	//	コンテナの終了
	if err = poolTest.Purge(resourceTest); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func connectDB() error {
	// DB(コンテナ)との接続
	if poolErr := poolTest.Retry(func() error {
		// DBコンテナが立ち上がってから疎通可能になるまで少しかかるのでちょっと待ったほうが良さそう
		time.Sleep(time.Second * 20)

		var err error
		connTest, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/GO_POSTGRES_GRAPHQL_DB?charset=utf8mb4&parseTime=True&loc=Local", resourceTest.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return connTest.Ping()
	}); poolErr != nil {
		log.Fatalf("Could not connect to docker: %s", poolErr)
		return poolErr
	}
	return nil
}

func execSQLScript(path string) error {
	var err error
	content, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return fileErr
	}
	_, err = connTest.Exec(bytes.NewBuffer(content).String())
	if err != nil {
		return err
	}
	return nil
}

func createTestData() error {
	var err error
	if err = testdata.CreateTestData(connTest); err != nil {
		return err
	}
	return nil
}
