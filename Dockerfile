# syntax=docker/dockerfile:1

FROM golang:1.24.1-alpine3.21 AS build

RUN wget https://musl.libc.org/releases/musl-1.2.5.tar.gz && \
   tar -xzf musl-1.2.5.tar.gz && \
   cd musl-1.2.5 && \
   ./configure --enable-static --disable-shared && \
   make && make install

RUN mkdir -pv /app
COPY . /app
WORKDIR /app

ENV GOPATH=/app
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN CGO_ENABLED=1 CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode=external -extldflags=-static' ./cmd/metagen
RUN ./cmd/metagen --env=production build

RUN CGO_ENABLED=1 CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode=external -extldflags=-static' ./cmd/server

FROM scratch

COPY --from=build /app/server /app/server
COPY --from=build /app/wwwroot /app/wwwroot
COPY --from=build /app/passwords.db /app/passwords.db

CMD ["/app/server"]