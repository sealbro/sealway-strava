# sealway-strava
sealway strava integration


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
- Replace `model.` to `strava.` without mutations
