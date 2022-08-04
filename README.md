# Prescaling-exporter

This project is a Prometheus exporter written in Golang, its goal is to provide a metric that scales applications to a requested capacity on a daily time slot or on an event declared in the embedded api. 

The project includes an exporter that calculates the metric to provide according to the annotations defined on the hpa. and an API that allows to declare an event and a capacity multiplier

## Requirements

   - Prometheus Stack or Victoria Metrics Stack
   - [Prometheus Adapater](https://github.com/kubernetes-sigs/prometheus-adapter) 

*Info : It is quite possible to use this solution with another observability stack than Prometheus. For example Datadog or Newrelic, but we don't have a configuration example.* 

# Install
## Kubernetes Deployement

- Clone repo
- Run this command with Helm3

```bash
helm install prescaling-exporter ./helm/prescaling-exporter -n prescaling-exporter --create-namespace
```

*You can use skaffold if you want.* 

## Configure a Horizontal Pod Autoscaler

To be able to pre-scale an application every day before a traffic spike, the only thing to do is to add the
following annotations on the HPA:

Annotations | values
--- | --- 
annotations.scaling.exporter.time.start | "hh:mm:ss"
annotations.scaling.exporter.time.end | "hh:mm:ss"
annotations.scaling.exporter.replica.min  | "integer"


```
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: "{{ .Release.Name }}"
  annotations:
    annotations.scaling.exporter.replica.min: "{{ .Values.hpa.annotations.replica_min"
    annotations.scaling.exporter.time.end: "{{ .Values.hpa.annotations.time_end }}"
    annotations.scaling.exporter.time.start: "{{ .Values.hpa.annotations.time_start }}"
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: "{{ .Release.Name }}"
  minReplicas: {{ .Values.hpa.minReplicas }}
  maxReplicas: {{ .Values.hpa.maxReplicas }}
  metrics:
  - type: External
    external:
      metricName: "min_replica"
      metricSelector:
          matchLabels:
            deployment: "{{ .Release.Name }}"
      targetValue: 10
```

*it's important to set the target to 10, the value provided for scale is 11, the addition of pod will be gradual. Currently, the increment is carried out 10% by 10%.*

## Configure prometheus adapter

A configuration example using the prometheus adapter to supply the metric to the Kubernetes cluster. 

```
  - "metricsQuery": "avg(<<.Series>>{<<.LabelMatchers>>})"
    "name":
      "as": "min_replica"
    "resources":
      "overrides":
        "namespace":
          "resource": "namespace"
    "seriesQuery": "min_replica"
```

## Configure parameters 

It is possible to change the application settings to use other annotations or the application port. This is done through configurable variable environments in the chart values

Parameters            | Default values                            | Comment
---                   | ---                                       | --- 
Namespace             | "prescaling-exporter"                     | Namespace for PrescalingEvent
Port                  | "9101"                                    | Application port
AnnotationEndTime     | "annotations.scaling.exporter.time.end"   | Annotation end in HPA for create metrique
AnnotationStartTime   | "annotations.scaling.exporter.time.start" | Annotation start in HPA for create metrique
AnnotationMinReplicas | "annotations.scaling.exporter.replica.min"| Annotation in min HPA for create metrique
LabelProject          | "project"                                 | label k8s where to retrieve the value for the project label of the metric

# Event Prescaling

With the CRDs and the API, it is possible to create prescaling events to allow the platform to scale on other schedules and with a multiplier.



# Build
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
