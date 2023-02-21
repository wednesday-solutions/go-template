FROM golang:1.18-alpine3.16 as builder
RUN apk add build-base

RUN mkdir  /app
ADD . /app

WORKDIR /app
ENV ENVIRONMENT_NAME=docker
RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor


RUN go run ./cmd/seeder/main.go
RUN go build -o ./output/server ./cmd/server/main.go
RUN go build -o ./output/migrations ./cmd/migrations/main.go
RUN go build  -o ./output/seeder ./cmd/seeder/exec/seed.go


FROM alpine:latest
RUN apk add --no-cache libc6-compat 
RUN apk add --no-cache --upgrade bash
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot


ARG ENVIRONMENT_NAME docker

RUN mkdir -p /app/
WORKDIR /app
USER nonroot

COPY /scripts /app/scripts/
COPY --from=builder /app/output/ /app/
COPY --from=builder /app/cmd/seeder/exec/build/ /app/cmd/seeder/exec/build/
COPY ./.env.docker /app/output/
COPY ./.env.docker /app/output/cmd/seeder/exec/build/
COPY ./.env.docker /app/output/cmd/seeder/exec/
COPY ./.env.docker /app/output/cmd/seeder/
COPY ./.env.docker /app/output/cmd/
COPY ./.env.docker /app/
COPY ./scripts/ /app/
COPY --from=builder /app/internal/migrations/ /app/internal/migrations/
CMD ["bash","./migrate-and-run.sh"]


EXPOSE 9000

