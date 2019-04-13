FROM golang:alpine AS build
ADD . /src

RUN apk add git

RUN cd /src && \
    go mod vendor && \
    go build -o ecs-go

# final stage
FROM alpine as bare
WORKDIR /app

COPY --from=build /src/ecs-go /app/

RUN apk add jq
ENV PATH="/app:${PATH}"
ENTRYPOINT ecs-go

