# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
WORKDIR /go/src/todolist
COPY . .

RUN bash get.sh

#RUN go get -u github.com/lib/pq
RUN go install github.com/yfedoruck/todolist

ENTRYPOINT /go/bin/todolist

# Document that the service listens on port 8080.
EXPOSE 8080
