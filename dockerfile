FROM golang:1.24-alpine AS build

WORKDIR /src/server
COPY server/ ./
RUN GO111MODULE=off CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app
COPY --from=build --chown=app:app /out/server /usr/local/bin/server

WORKDIR /app
COPY --chown=app:app html/ ./html/
COPY --chown=app:app js/ ./js/
COPY --chown=app:app photo/ ./photo/

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]