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
  
## Алгоритм работы и возможности

- API Получение ИД события от стравы
  - Положить в Очередь или локальное хранилище в случае проблемы
  - Положить в БД
  - Создать действие я получение данных или обновление
- Сервис/демон который учитывая количество запросов ставит в очередь загрузки данные
  - сохранение в БД
- gRPC/API получает запрос по пользователям (AthleteId, token, priority) с приоритетами для загрузки данных
- если данных нет то возвращается список AthleteId + NextRequestDate + WaitActivities(ExpectedDate, ActivityId)
- если данные есть и готовы, то возвращаются или ActivityId[] c которыми нужно будет идти в GraphQL
