FROM golang

RUN mkdir -p /go/src/github.com/wednesday-solutions/go-template

RUN go install github.com/rubenv/sql-migrate/... \
  github.com/volatiletech/sqlboiler \
  github.com/99designs/gqlgen

ADD . /go/src/github.com/wednesday-solutions/go-template
WORKDIR /go/src/github.com/wednesday-solutions/go-template

RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor
RUN go build -o ./ ./cmd/server/main.go
CMD ["bash", "./migrate-and-run.sh"]
EXPOSE 9000