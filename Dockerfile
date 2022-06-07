FROM golang:1.18-alpine as build

WORKDIR /app

COPY src/go.mod ./
COPY src/go.sum ./

RUN go mod download

COPY src/*.go ./

RUN go build -o /opencollective-webhook

FROM alpine:3.14.1

WORKDIR /usr/app

COPY --from=build /opencollective-webhook /opencollective-webhook

ENTRYPOINT ["/opencollective-webhook"]