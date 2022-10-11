# first stage - builds the binary from sources
FROM golang:alpine as build

# using build as current directory
WORKDIR /app

# adding the source code to current dir:
COPY . .

# downloading dependencies and verifying
RUN go mod download && go mod verify


# building the project
RUN go build  -o food-ordering .

# second stage - using minimal image to run the server
FROM alpine:latest

# using /app as current directory
WORKDIR /app

# copy server binary from `build` layer
COPY --from=build /app/food-ordering .


# binary to run
CMD "/app/food-ordering"

EXPOSE 8084