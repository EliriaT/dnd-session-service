# first stage - builds the binary from sources
FROM golang:1.23.9-alpine as build

# using build as current directoryAdd commentMore actions
WORKDIR /app

# adding the source code to current dir:
COPY . .

# downloading dependencies and verifying
RUN go mod download && go mod verify

# building the project
RUN go build  -o main main.go

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# second stage - using minimal image to run the server
FROM alpine:latest

# using /app as current directory
WORKDIR /app

# copy server binary from `build` layer
COPY --from=build /app/main .
COPY --from=build /app/migrate ./migrate
COPY --from=build /app/app.env ./app.env
COPY wait-for.sh .
COPY start.sh .
COPY db/migrations ./migrations

# binary to run

CMD ["/app/main"]

ENTRYPOINT [ "/app/start.sh" ]

EXPOSE 8082