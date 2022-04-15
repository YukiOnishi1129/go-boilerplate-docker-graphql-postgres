# go-boilerplate-docker-graphql-postgres

golang GraphQL ボイラーテンプレート

## 技術構成
- go: 1.18
- postgres
- gqlgen
- sqlboiler
- jwt-go
- aws-sdk-go
- ozzo-validation
- staticcheck
- golang-migrate
- github actions

## 仕様
- ログイン機能付きのTodoリスト
  - 認証処理はcookie形式を採用 (userIDをjwtトークンに変換し、http onlyのcookieにセット)

## 主な機能
### 認証権限なし
- ログイン
- 会員登録
- ログアウト

### 認証権限あり
- Todo
  - Todoリスト取得
  - Todo詳細取得
  - Todo新規登録
  - Todo更新処理
  - Todo削除処理
- User
  - ユーザー情報取得
  - ユーザー名変更処理
  - メールアドレス変更処理
  - パスワード変更処理
  - ユーザー画像登録、変更処理


## GRAPHQLスキーマ

### User
- https://github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/blob/main/graphql/user.graphql
### Todo
- https://github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/blob/main/graphql/todo.graphql


## 環境構築
### 1. envファイルを作成
```
// ルートディレクトリ直下に「.env」ファイルを作成
touch .env

//「.env.sample」の記述を「.env」にコピー
// ※AWS関連は各自用意

// appディレクトリ直下に「.env」ファイルを作成
touch app/.env

//「app/.env.sample」の記述を「app/.env」にコピー
```

### 2. docker起動
- ビルド
```
docker-compose build
```

- コンテナ起動
```
docker-compose up -d
```

### 3. マイグレーション
- golang-migrateをmacにインストール
```
 brew install golang-migrate
```

- マイグレーションを実行
- ※ DBコンテナを事前に起動しておくこと
```
// ルートディレクトリで実行
./dev-tools/bin/db:migrate
```

- シーダー実行
```
// ルートディレクトリで実行
./dev-tools/bin/db:seed
```

### 4. GraphQL実行環境
- `localhost:4000`でGraphiQLが立ち上がる
  - graphqlのクエリーを実行できるようになります。
  - ![スクリーンショット 2022-04-16 6 49 13](https://user-images.githubusercontent.com/58220747/163649363-ae18280f-cab9-42f1-91aa-6ac1637ebc44.png)

  
  - クエリについては以下を参照
- query, mutation, fragmentは以下を記載
```
# Todoリスト一覧取得
query getTodoList {
  todoList {
    ...todoDetail
  }
}

# 単一のTodoを取得
query getTodoDetail($todoId: ID!) {
  todoDetail(id:$todoId) {
    ...todoDetail
  }
}

# Todo新規作成
mutation createTodoDetail($createInput: CreateTodoInput!) {
  createTodo(input: $createInput){
    ...todoDetail
  }
}

# Todo更新
mutation updateTodoDetail($updateInput: UpdateTodoInput!) {
  updateTodo(input: $updateInput){
    ...todoDetail
  }
}

# Todo削除
mutation deleteTodoDetail($deleteId: ID!) {
  deleteTodo(id: $deleteId)
}


fragment todoDetail on Todo {
  id
  title
  comment
  user {
    ...userDetail
  }
  createdAt
  updatedAt
  deletedAt
}

# ユーザー情報取得
query getUserDetail($userId: ID!) {
  userDetail(id: $userId) {
    ...userDetail
  }
}

# ログイン
mutation SignInDetail($signinInput: SignInInput!) {
  signIn(input: $signinInput) {
    ...userDetail
  }
}

#会員登録
mutation SignUpDetail($signupInput: SignUpInput!) {
  signUp(input: $signupInput) {
    ...userDetail
  }
}

# ログアウト
mutation SignOutDetail {
  signOut
}

# ユーザー名変更処理
mutation UpdateUserNameDetail($userName: String!) {
  updateUserName(name: $userName) {
    ...userDetail
  }
}

# メールアドレス変更処理
mutation UpdateUserEmailDetail($userEmail: String!) {
  updateUserEmail(email: $userEmail) {
    ...userDetail
  }
}

# パスワード変更処理
mutation UpdateUserPasswordDetail($userPassword: updatePasswordInput!) {
  updateUserPassword(input: $userPassword) {
    ...userDetail
  }
}

# ユーザー画像変更処理
# AWS S3の設定と、Altair GraphQL Clientでの実行確認が必要です。
mutation UploadUserFileDetail($userFile:Upload!) {
  uploadUserFile(file: $userFile) {
    ...userDetail
  }
}

fragment userDetail on User {
  id
  name
  email
  imageUrl
  createdAt
  updatedAt
  deletedAt
}
```

- query variablesは以下を記載
```
{
  "todoId": "1",
  "createInput": {
    "title": "サンプル",
    "comment": "サンプル"
  },
  "updateInput": {
    "id": "2",
     "title": "サンプル111",
    "comment": "サンプル111"
  },
  "deleteId": "4",
  "signupInput": {
    "name": "たかし",
    "email": "taro@gmail.com",
    "password": "password",
    "passwordConfirm": "password"
  },
  "signinInput": {
    "email": "taro@gmail.com",
    "password": "passwd"
  },
  "userId": "2",
  "userName": "透",
  "userEmail": "taro_ver2@gmail.com",
  "userPassword": {
    "oldPassword": "password",
    "newPassword": "passwd",
    "newPasswordConfirm": "passwd"
  }
}
```

- ユーザー画像作成・変更処理はAWS S3の設定と、Altair GraphQL Clientでの実行確認が必須です。
  - Alter GraphQL Clientを用いた画像アップロードの確認方法についてはこちら
  - https://www.wantedly.com/companies/visitsworks/post_articles/330336

## 開発用コマンド
- 以下のコマンドは全てルートディレクトリで実行すること
### モデルファイル自動生成
  - SqlBoilerを用いて、テーブル構造からモデルを生成する (gormとは逆のパターンで生成)
  - ※事前にDBコンテナを立ち上げておくこと
```
./dev-tools/bin/runner.sh entity:create
```

### マイグレーション
- テーブル作成
- ※事前にDBコンテナを立ち上げておくこと
```
./dev-tools/bin/runner.sh db:migrate
```

### シーディング
- データ登録
- ※事前にDBコンテナを立ち上げておくこと
```
./dev-tools/bin/runner.sh db:seed
```

### ロールバック
- テーブル設定を初期化
- ※事前にDBコンテナを立ち上げておくこと
```
./dev-tools/bin/runner.sh db:rollback
```

### データ初期化
- ロールバック、マイグレーション、シーディングを一気に実行してテーブルデータを初期化する
- ※事前にDBコンテナを立ち上げておくこと
```
./dev-tools/bin/runner.sh db:reset
```

### 静的解析 (lint)
```
./dev-tools/bin/runner.sh lint
```

### テスト
全てのテストを実行
```
./dev-tools/bin/runner.sh test:all
```

### graphql generate
- graphqlスキーマファイルからresolverなどを自動生成
```
./dev-tools/bin/runner.sh gql
```
