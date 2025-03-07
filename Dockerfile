# syntax=docker/dockerfile:1

FROM alpine:3.14
RUN apk add --no-cache git make musl-dev go
