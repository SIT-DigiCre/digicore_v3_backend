# digicore v3 backend

## setup env

1. .env.sampleをコピーして.envを作成する
1. [Discord developers](https://discord.com/developers/applications)でAppを作成し、Oauth2のRedirectsに`${BACKEND_ROOT_URL}/discord/oauth/callback`を指定する
1. 上記で作成したAppのClient informationからClient IDとClient Secretを取得し、.envに追記する。
1. [Google Cloud Platform](https://console.cloud.google.com/home/dashboard)でAppを作成し、OAuth クライアント IDをアプリケーションの種類をウェブアプリケーションにして作成し、承認済みのリダイレクト URIに`${BACKEND_ROOT_URL}/google/oauth/callback/login`と`${BACKEND_ROOT_URL}/google/oauth/callback/register`を指定する。
1. 上記で作成したAppのclient_secret_*.jsonをダウンロードし、client_secret.jsonに名前を書き換えこのファイルが有る階層に配置する。
1. [build env](#build-env)を行う
1. [run env](#run-env)を行う
1. [db migration](#db-migration)を行う

## build env

```sh
docker compose build
```

## run env

```sh
docker compose up
```

## db migration

```sh
docker compose run -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.11.50 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --dry-run'
docker compose run -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.11.50 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE}'
docker compose run -w /app/schema admin bash -c 'mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} < foreign_key.sql'
```

## generate swagger docs

```sh
swag init
```
