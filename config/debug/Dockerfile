# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine AS build-env
RUN apk add --no-cache git

ENV CGO_ENABLED 0
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/yfedoruck/todolist

RUN go get -u github.com/lib/pq

RUN go build -gcflags "all=-N -l" -o /server /go/src/github.com/yfedoruck/todolist

# Compile Delve
RUN go get github.com/derekparker/delve/cmd/dlv

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install github.com/yfedoruck/todolist

# Run the outyet command by default when the container starts.
#ENTRYPOINT /go/bin/todolist

# Final stage
FROM alpine:3.7
# Document that the service listens on port 8080.
EXPOSE 8080 40000

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=build-env /server /
COPY --from=build-env /go/bin/dlv /

# Run delve
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/server"]
