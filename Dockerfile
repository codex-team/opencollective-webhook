FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o /opencollective-webhook

FROM alpine:3.14.1

WORKDIR /usr/app

COPY --from=build /opencollective-webhook /opencollective-webhook

COPY .env.sample ./

ENTRYPOINT ["/opencollective-webhook"]