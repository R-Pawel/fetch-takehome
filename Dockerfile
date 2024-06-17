FROM golang:1.22-alpine as builder

WORKDIR /project/go-docker/


ENV GIN_MODE=release

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /project/go-docker/build/myapp .

FROM alpine:latest
COPY --from=builder /project/go-docker/build/myapp /project/go-docker/build/myapp

EXPOSE 8080
ENTRYPOINT [ "/project/go-docker/build/myapp" ]