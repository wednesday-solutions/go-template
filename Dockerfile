FROM golang:1.18-alpine3.16
RUN apk add build-base

RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template

ADD . /go/src/github.com/wednesday-solutions/go-template

WORKDIR /go/src/github.com/wednesday-solutions/go-template

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template/output
RUN go build -o ./output/main ./cmd/server/main.go


FROM golang:1.18-alpine3.16
ARG ENVIRONMENT_NAME
RUN apk add build-base
RUN mkdir -p /app

ADD .  /app

WORKDIR /app

COPY --from=0 /go/src/github.com/wednesday-solutions/go-template/output /app/

CMD ["sh", "/app/scripts/migrate-and-run.sh"]
EXPOSE 9000