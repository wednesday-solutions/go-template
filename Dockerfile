FROM golang:1.22-alpine3.19 AS builder
RUN apk add build-base

ARG APP_PATH

RUN mkdir -p ${APP_PATH}
ADD . ${APP_PATH}

WORKDIR ${APP_PATH}
ARG ENVIRONMENT_NAME 
ENV ENVIRONMENT_NAME=$ENVIRONMENT_NAME
RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go mod vendor


RUN go run ./cmd/seeder/main.go
RUN go build -cover -o ./output/server ./cmd/server/main.go
RUN go build -cover -o ./output/migrations ./cmd/migrations/main.go
RUN go build -cover -o ./output/seeder ./cmd/seeder/exec/seed.go


FROM alpine:latest
RUN apk add --no-cache libc6-compat 
RUN apk add --no-cache --upgrade bash
RUN apk add --no-cache dumb-init
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

ARG APP_PATH
ARG ENVIRONMENT_NAME
ENV ENVIRONMENT_NAME=$ENVIRONMENT_NAME

RUN echo ${APP_PATH}
RUN mkdir -p ${APP_PATH}/
WORKDIR ${APP_PATH}

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Clean the coverage-reports folder
USER nonroot

COPY /scripts ${APP_PATH}/scripts/
COPY --from=builder ${APP_PATH}/output/ ${APP_PATH}/
COPY --from=builder ${APP_PATH}/cmd/seeder/exec/build/ ${APP_PATH}/cmd/seeder/exec/build/

COPY ./.env.* ${APP_PATH}/output/
COPY ./.env.* ${APP_PATH}/output/cmd/seeder/exec/build/
COPY ./.env.* ${APP_PATH}/output/cmd/seeder/exec/
COPY ./.env.* ${APP_PATH}/output/cmd/seeder/
COPY ./.env.* ${APP_PATH}/output/cmd/
COPY ./.env.* ${APP_PATH}/
COPY ./scripts/ ${APP_PATH}/
COPY --from=builder ${APP_PATH}/internal/migrations/ ${APP_PATH}/internal/migrations/

ENTRYPOINT ["dumb-init", "--"]
STOPSIGNAL SIGINT
CMD ["/entrypoint.sh"]
EXPOSE 9000

