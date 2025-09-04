FROM golang:1.22.3-alpine3.19

ENV repo /go/src/github.com/Orange-Health/citadel

WORKDIR ${repo}

RUN go install github.com/pilu/fresh@latest

COPY go.mod ${repo}
COPY go.sum ${repo}
RUN go mod download

ADD . ${repo}

CMD fresh -c ${repo}/fresh.conf;

EXPOSE 8080
