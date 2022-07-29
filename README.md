# Prescaling-exporter

This project is a Prometheus exporter written in Golang, its goal is to provide a metric that scales applications to a requested capacity on a daily time slot or on an event declared in the embedded api. 

The project includes an exporter that calculates the metric to provide according to the annotations defined on the hpa. and an API that allows to declare an event and a capacity multiplier

## Requirements

   - Prometheus Stack or Victoria Metrics Stack
   - [Prometheus Adapater](https://github.com/kubernetes-sigs/prometheus-adapter) 

## Install

- Clone repo
- Run this command with Helm3

```bash
helm install prescaling-exporter ./kubernetes/helm/prescaling-exporter
```

## Configure a Horizontal Pod Autoscaler

To be able to pre-scale an application every day before a traffic spike, the only thing to do is to add the
following annotations on the HPA:

Annotations | values
--- | --- 
annotations.scaling.exporter.time.start | "hh:mm:ss"
annotations.scaling.exporter.time.end | "hh:mm:ss"
annotations.scaling.exporter.replica.min  | "integer"


## Golang utils 

1. Golang build 

```bash
./generate_type.sh
go build
```

2. Run test

```bash
go test -v ./...
```
3. Run test with coverate report

```bash
 go test -coverprofile=cover.out -v ./...   
```

```
openapi-generator generate -i api/openapi.yaml -g go-server -o .
```

## Generate crd files

```
./generate_type.sh
```

## Generate swagger docs

```
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g pkg/server/server.go --parseInternal=true
```
