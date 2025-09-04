FROM golang:1.22.3-alpine3.19

ENV repo /go/src/github.com/Orange-Health/citadel

WORKDIR ${repo}

# RUN go get github.com/pilu/fresh

COPY go.mod ${repo}
COPY go.sum ${repo}
RUN go mod download

ADD . ${repo}

# CMD fresh -c ${repo}/worker/fresh.conf;

RUN go build -o /go/bin/consumer ${repo}/main/consumer

ENTRYPOINT [ "/go/bin/consumer" ]
