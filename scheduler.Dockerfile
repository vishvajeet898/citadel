FROM golang:1.22.3-bullseye

ENV repo /go/src/github.com/Orange-Health/citadel

WORKDIR ${repo}

COPY go.mod ${repo}
COPY go.sum ${repo}
RUN go mod download

ADD . ${repo}

RUN apt-get update -y && apt-get install -y s4cmd jq awscli curl

RUN chmod +x ${repo}/docker-entrypoint.sh

RUN go build -o /go/bin/scheduler ${repo}/main/scheduler

ENTRYPOINT [ "./docker-entrypoint.sh" ]

CMD [ "/go/bin/scheduler" ]
