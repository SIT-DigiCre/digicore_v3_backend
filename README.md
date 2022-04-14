# digicore v3 backend

## setup env

1. .env.sampleをコピーして.envを作成する
1. [Discord developers](https://discord.com/developers/applications)でAppを作成し、Oauth2のRedirectsに`http://localhost:8000/discord/oauth/callback`を指定する
1. 上記で作成したAppのClient informationからClient IDとClient Secretを取得し、.envに追記する。
1. [Google Cloud Platform](https://console.cloud.google.com/home/dashboard)でAppを作成し、OAuth クライアント IDをアプリケーションの種類をウェブアプリケーションにして作成し、承認済みのリダイレクト URIに`http://localhost:8000/google/oauth/callback`を指定する。
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
docker compose exec -w /app backend  bash -c 'cat ./schema/*.sql | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.11.50 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --dry-run'
docker compose exec -w /app backend  bash -c 'cat ./schema/*.sql | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.11.50 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE}'
```

## generate swagger docs

```sh
swag init
```
