# WanRun

## docker command

### 1. postgresのマウント用のフォルダ作成

```
mkdir -p /var/postgres
```
### 2.1 ホットリロードで開発するのか実行環境で立ち上げるのか
`docker-compose.yaml`の`backendコンテナ`を下記に設定。
**ホットリロードバージョン**
```
target: Dev
```

**実行環境で立ち上げたい場合**
```
target: Deploy
```

### 2.2 Postgres立ち上げとEchoの立ち上げ

```
docker compose up -d --build
```

### 2.3 コンテナ２台の確認(backend, postgres)

```
docker ps
```

### 3. FYI

https://github.com/air-verse/air

---

## Local db migration

### FYI
https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

### SQLファイル作成コマンド
```
*絶対パスで指定ができなかったので注意(ルートフォルダからやって)
migrate create -ext sql -dir migrate/migration_sql -seq {create_table_name}
```

### DB作成コマンド
dbディレクトリの中に移動して下記のコマンドを実行
`MIGRATION_DIRECTION`の環境変数で識別をしてる。

**DBテーブル作成**
```
MIGRATION_DIRECTION=up go run migrate.go
```

**DBテーブル削除**
```
MIGRATION_DIRECTION=down go run migrate.go
```

**DB削除**
※Dirty database versionエラーが発生した場合
マイグレーションがどのバージョンで失敗したかを確認します。これはデータベースのマイグレーション履歴テーブル（通常はschema_migrationsまたはflyway_schema_history）を直接クエリすることで確認
```
MIGRATION_DIRECTION=force MIGRATION_VERSION={対象のバージョン番号} go run migrate.go
```

## go doc

### インストール方法
```
go install golang.org/x/tools/cmd/godoc@latest
```

### 閲覧方法
下記のコマンド後に、`http://localhost:6060`を開く
```
godoc -http=:6060
```

### internalなどの部分の閲覧方法
`/cmd`, `/internal`, `/pkg`
などは、下記のパスを参考にすること

> As mentioned in this documentation, using m=all parameter in URL will list internal packages.


ex) 
`http://localhost:6060/pkg/github.com/wanrun-develop/wanrun/internal/auth/core/handler/?m=all`

FYI: https://stackoverflow.com/questions/67185294/cant-godoc-create-documentation-for-packages-within-internal-folder

## Folder Hierarchy
https://github.com/golang-standards/project-layout

## Naming convention
https://go.dev/doc/effective_go#names
