FROM golang:alpine AS build
ADD . /src

RUN apk add -U --no-cache ca-certificates git make

RUN cd /src && \
    make install && \
    go build -o ecs-go

# final stage
FROM alpine as bare
WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/ecs-go /app/

RUN apk add jq
ENV PATH="/app:${PATH}"
ENTRYPOINT ecs-go

