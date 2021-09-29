FROM golang:alpine as builder

WORKDIR /build
COPY . .

RUN go build

FROM alpine as runtime

WORKDIR /app

COPY --from=builder /build/sealway-strava .

ENTRYPOINT ./sealway-strava