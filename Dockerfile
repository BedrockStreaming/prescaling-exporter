FROM golang:1.18

ENV GO111MODULE=off

RUN go get k8s.io/code-generator k8s.io/apimachinery

ARG repo="${GOPATH}/src/github.com/bedrockstreaming/prescaling-exporter"

RUN mkdir -p $repo
WORKDIR $repo
VOLUME $repo
