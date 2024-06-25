FROM golang:1.22-alpine3.19 AS builder
RUN apk add build-base

WORKDIR /app
# copy the go.mod and go.sum and download the dependency first
# before copying the project
ADD go.mod /app
ADD go.sum /app
RUN go mod download

# NOW ADD the whole root project
ADD . /app

ARG ENVIRONMENT_NAME 
ENV ENVIRONMENT_NAME=$ENVIRONMENT_NAME
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


ARG ENVIRONMENT_NAME
ENV ENVIRONMENT_NAME=$ENVIRONMENT_NAME

RUN mkdir -p /app/
WORKDIR /app
USER nonroot

COPY /scripts /app/scripts/
COPY --from=builder /app/output/ /app/
COPY --from=builder /app/cmd/seeder/exec/build/ /app/cmd/seeder/exec/build/

COPY ./.env.* /app/output/
COPY ./.env.* /app/output/cmd/seeder/exec/build/
COPY ./.env.* /app/output/cmd/seeder/exec/
COPY ./.env.* /app/output/cmd/seeder/
COPY ./.env.* /app/output/cmd/
COPY ./.env.* /app/
COPY ./scripts/ /app/
COPY --from=builder /app/internal/migrations/ /app/internal/migrations/
CMD ["bash","./migrate-and-run.sh"]


EXPOSE 9000

