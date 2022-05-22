FROM golang:alpine as builder

WORKDIR /build
COPY . .

RUN go build

#FROM alpine as runtime
FROM gcr.io/distroless/static as runtime

COPY --from=builder /build/sealway-strava .

EXPOSE 8080

ENTRYPOINT ./sealway-strava
