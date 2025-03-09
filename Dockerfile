# syntax=docker/dockerfile:1

FROM golang:1.24.1-alpine3.21 AS build

RUN apk add clang lld musl-dev compiler-rt compiler-rt-static

RUN mkdir -pv /app
COPY . /app
WORKDIR /app

ENV CC=clang
ENV GOPATH=/app
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN go build --ldflags '-linkmode=external -extldflags=-static' ./cmd/metagen
RUN ./cmd/metagen --env=production build

RUN go build --ldflags '-linkmode=external -extldflags=-static' ./cmd/server

FROM scratch

COPY --from=build /app/server /app/server
COPY --from=build /app/wwwroot /app/wwwroot
COPY --from=build /app/passwords.db /app/passwords.db

CMD ["/app/server"]