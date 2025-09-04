FROM golang:1.22.3-alpine3.19

ENV repo /go/src/github.com/Orange-Health/citadel

#Cgo enables the creation of Go packages that call C code.
ENV CGO_ENABLED 0

WORKDIR ${repo}

COPY go.mod ${repo}
COPY go.sum ${repo}
RUN go mod download

# The inotify API provides a mechanism for monitoring filesystem
# events. Inotify can be used to monitor individual files, or to
# monitor directories. When a directory is monitored, inotify will
# return events for the directory itself, and for files inside the
# directory.
RUN apk update && apk add bash inotify-tools git

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ADD . ${repo}

EXPOSE 8080
RUN go build -o app

### Run the Delve debugger ###
COPY ./startScript.sh /
RUN chmod +x /startScript.sh
ENTRYPOINT [ "/startScript.sh"]
