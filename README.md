# go-chat-ai-server

OpenAI Chat API test

Please prepare the environment for postgres to use as the DB.

## setup

```shell
% cp .envrc.example .envrc
```

### edit .envrc
``` .envrc
export OPEN_API_TOKEN=[your open ai api token]
export MY_API_TOKEN=[your api call api token]
export PORT=[your port]
export DB_USER=[your postgres user]
export DB_PASSWORD=[your postgres password]
export DB_HOST=[your postgres host]

```

## usage

``` shell
% go run main.go
```

```shell
% curl 'http://localhost:[your port]/api?text=こんにちは' -H 'Content-Type: application/json' -H 'Authorization: Bearer [your api call api token]'

# example
% curl 'http://localhost:8080?/api?text=こんにちは' -H 'Content-Type: application/json' -H 'Authorization: Bearer api_key_1234' 
```

※本アプリケーションを利用する場合、簡易的な認証機構を入れていますが広く公開する場合はきちんと認証機構を実装することをお勧めします