### Go get dependecies and build ###
FROM golang:1.16 as builder
WORKDIR /go/src/app
COPY go.mod go.sum ./

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly

### Create image ###
FROM alpine:3
WORKDIR /usr/bin
COPY --from=builder /go/src/app/oasdiff .
ENTRYPOINT ["/usr/bin/oasdiff"]
