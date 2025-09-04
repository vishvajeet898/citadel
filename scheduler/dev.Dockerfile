FROM golang:1.22.3-alpine3.19

ENV repo /go/src/github.com/Orange-Health/citadel

WORKDIR ${repo}

COPY go.mod ${repo}
COPY go.sum ${repo}
RUN go mod download

ADD . ${repo}

RUN go build -o /go/bin/scheduler ${repo}/main/scheduler

ENTRYPOINT [ "/go/bin/scheduler" ]
