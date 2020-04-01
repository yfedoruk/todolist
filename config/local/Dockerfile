###
 # docker exec todolist_server_1 go install && docker stop todolist_server_1 && docker start todolist_server_1
 # //or
 # docker exec todolist_server_1 go install &&  docker-compose.exe restart server

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# Copy the local package files to the container's workspace.
WORKDIR /go/src/todolist

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/todolist github.com/yfedoruck/todolist/cmd

ENTRYPOINT /go/bin/todolist

# Document that the service listens on port 8080.
EXPOSE 8080
