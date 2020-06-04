# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
# Copy the local package files to the container's workspace.
ADD . docker -ps
# Document that the service listens on port 8080.
EXPOSE 4242