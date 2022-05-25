# sealway-strava

Caching request service to strava.
Entities cached/stored in mongodb and you can get them with graphQL

![build](https://github.com/sealbro/sealway-strava/actions/workflows/docker.yml/badge.svg)
[![Docker Pulls](https://badgen.net/docker/pulls/sealway/strava?icon=docker&label=pulls)](https://hub.docker.com/r/sealway/strava/)
[![Docker Image Size](https://badgen.net/docker/size/sealway/strava?icon=docker&label=image%20size)](https://hub.docker.com/r/sealway/strava/)


## Environments

- `MONGO_CONNECTION` - mongo connection string
- `STRAVA_CLIENT` - strava client id
- `STRAVA_SECRET` - strava secret id
- `ACTIVITY_BATCH_SIZE` (50) - max batch size, after which data is sent to subscribers to subscribers
- `ACTIVITY_BATCH_TIME` (45s) - time after which data is sent to subscribers
- `SLUG` (integration-strava) - prefix for url path
- `PORT` (8080) - server port

## API Endpoints

- GET `{SLUG}/healthz` - health check
- GET `{SLUG}/api/quota` - actual strava's request quota
- GET `{SLUG}/api/subscription` - used for registration strava callback
- POST `{SLUG}/api/subscription` - there strava sends user changes
- GET `{SLUG}/graphql/` - graphQL playgroud
- [More](./interfaces/graph/schema.graphqls) about graphQL queries / mutations / subscriptions

## Debug

- set environments
  - `SEALWAY_ConnectionStrings__Mongo__Connection`
  - `SEALWAY_Services__Strava__Client`
  - `SEALWAY_Services__Strava__Secret`
- localhost mongo `docker run -d --restart=always --name mongodb -p 27017:27017 mongo`
- [graphql queries docs](https://graphql.org/learn/queries/)
- [timestamp converter](https://www.epochconverter.com/)

## Generate strava client

- ```git clone https://github.com/swagger-api/swagger-codegen.git```
- ```cd ./swagger-codegen```
  - for Windows in `./run-in-docker.sh` add `MSYS_NO_PATHCONV=1` before `docker run ...`
- ```./run-in-docker.sh generate --input-spec https://developers.strava.com/swagger/swagger.json --lang go --output /gen/go/```
- ```./run-in-docker.sh generate --input-spec https://developers.strava.com/swagger/swagger.json --lang openapi --output /gen/openapi/```
- Replace
    - ```package swagger``` to ```package strava```
    - ```(json:"([^,]+),)``` to ```bson:"$2" $1```
- gRPC [support](https://ednsquare.com/story/build-simple-api-with-grpc-protobuf-and-golang------kuxI0H)
    - ```go get -u github.com/golang/protobuf/protoc-gen-go```
    - ```protoc --go_out=plugins=grpc:. api.proto```

## Generate grahpql

- `cd ./graph`
- `go run github.com/99designs/gqlgen generate`
- Replace `model.` to `strava.` without mutations in `graph/generated/generated.go` and rollback `schema.resolvers.go`
  - `strava.NewAthleteToken` => `model.NewAthleteToken`
  - `strava.AthleteToken` => `model.AthleteToken`

## Ideas

- mutation reload all activities by athlete with limit
- mutation download activities from strava by athlete with limit

## Upgrade

- change go version in `go.mod`
- `go get -u all`

## Problems

- graphql generate `[Type!]!` or `[Type!]` [here](https://github.com/graph-gophers/graphql-go/issues/78#issue-220709670)
