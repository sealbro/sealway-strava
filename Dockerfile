FROM golang:1.18 as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/sealway-strava

FROM gcr.io/distroless/static as runtime

COPY --from=builder /bin/sealway-strava /

CMD ["/sealway-strava"]
