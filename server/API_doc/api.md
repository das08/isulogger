# ISULOGGER API仕様書

## コンテストの作成
新規コンテストを作成する。
- Endpoint  
[POST /contest](#post-contest)
```http request
http://localhost:8082/contest
```

- Request
```json
{
  "contest_name": "test contest"
}
```

- Response
  - Success: 200 OK, `contest_id`が返される
  - Error: 500 Internal Server Error

## ログエントリの作成
指定した`contest_id`のログエントリを作成する。(ログファイルはこのエンドポイントでは受け付けない)
- Endpoint  
  [POST /entry](#post-entry)
```http request
http://localhost:8082/entry
```

- Request
```json
{
  "contest_id": 3,
  "branch_name": "master",
  "score": 1000,
  "message": "initial bench"
}
```

- Response
    - Success: 200 OK, エントリの`id`が返される
    - Error: 500 Internal Server Error, 400 Bad Request

## ログエントリにログファイルの追加
指定した`contest_id`のログエントリにログファイルを追加する。追加方法は2種類存在する。  
- 最新のエントリにログファイルを追加する
- 指定したエントリ`id`にログファイルを追加する

### 最新のエントリにログファイルを追加する
指定した`contest_id`に一致する最新のエントリにログファイルを追加する。

- Endpoint  
  [POST /entry/:contest_id/:log_type](#post-logs)
```http request
http://localhost:8082/entry/:contest_id/:log_type
```

`log_type`は`access`か`slow`を指定する。

- Request
FORMでファイルをアップロードする。  
`log=<ログファイル>`


- Response
    - Success: 200 OK, エントリの`id`が返される
    - Error: 500 Internal Server Error, 400 Bad Request


### 指定したエントリIDにログファイルを追加する
指定したエントリ`id`に一致する最新のエントリにログファイルを追加する。

- Endpoint  
  [POST /entry/:contest_id/:log_type](#post-logs)
```http request
http://localhost:8082/entry/:contest_id/:log_type
```

`log_type`は`access`か`slow`を指定する。

- Request
  FORMでファイルと`id`をアップロードする。  
`log=<ログファイル>`  
`entry_id=<エントリID>`


- Response
    - Success: 200 OK, エントリの`id`が返される
    - Error: 500 Internal Server Error, 400 Bad Request

