FROM golang:1.16-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /opencollective-webhook

FROM alpine:3.14.1

WORKDIR /

COPY --from=build /opencollective-webhook /opencollective-webhook

COPY .env ./

ENTRYPOINT ["/opencollective-webhook"]