# Elon The Great Service

- Http service that return statistics about Elon musk tweets based on imported file


### Runing:
- Building and running:
    `go build cmd/tweet/main.go` or
    `go run cmd/tweet/main.go`
- Docker:
    `docker-compose build`
    `docker-compose up`

You can configure:
- PORT on which the server is using by setting PORT environment variable. Default value is 9000
- MONGO_URI: uri for mongodb client to connect
- TWEET_FILE: file path with tweets to be imported in db
- INIT_SEED: if included this will start importer for TWEET_FILE before running the service