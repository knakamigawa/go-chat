# go-chat-ai-server

OpenAI Chat API test

## setup

```shell
% cp .envrc.example .envrc
```

### edit .envrc
``` .envrc
export OPEN_API_TOKEN=[your open ai api token]
export MY_API_TOKEN=[your api call api token]
export PORT=[your port]
```

## usage

``` shell
% go run main.go
```

```shell
% curl 'http://localhost:[your port]?text=こんにちは' -H 'Content-Type: application/json' -H 'Authorization: Bearer [your api call api token]'

# example
% curl 'http://localhost:8080?text=こんにちは' -H 'Content-Type: application/json' -H 'Authorization: Bearer api_key_1234' 
```