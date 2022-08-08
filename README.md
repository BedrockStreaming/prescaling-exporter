# Prescaling-exporter

This project is a Prometheus exporter written in Golang, its goal is to provide a metric that scales applications to a requested capacity on a daily time slot or on an event declared in the embedded api. 

The exporter exposes a metric calculated according to the annotations defined on the HPA. The project also exposes an API that allows to register an event with a capacity multiplier.

## Requirements

   - Prometheus Stack or Victoria Metrics Stack
   - [Prometheus Adapater](https://github.com/kubernetes-sigs/prometheus-adapter) 

> Info: It is quite possible to use this solution with another observability stack than Prometheus. For example, Datadog or Newrelic, but we do not provide a configuration example.

# Install
## Kubernetes Deployement

- Clone repo
- Run this command with Helm3

```bash
helm install prescaling-exporter ./helm/prescaling-exporter -n prescaling-exporter --create-namespace
```

> You can use skaffold if you want.

## Configure an Horizontal Pod Autoscaler

To be able to pre-scale an application every day before a traffic spike, you must add the
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
      metricName: "prescale_metric"
      metricSelector:
          matchLabels:
            deployment: "{{ .Release.Name }}"
      targetValue: 10
```

> It's important to set the `targetValue` to 10. The metric's value provided by the exporter in order to scale is 11. The scale up of pods will be gradual. Currently, the increment is carried out 10% by 10%.

## Configure prometheus adapter

Here is a configuration example using the prometheus adapter to supply the metric to the Kubernetes cluster:

```
  - "metricsQuery": "avg(<<.Series>>{<<.LabelMatchers>>})"
    "name":
      "as": "prescale_metric"
    "resources":
      "overrides":
        "namespace":
          "resource": "namespace"
    "seriesQuery": "prescale_metric"
```

## Configure parameters 

You can modify the application's settings to use different annotations or a different port. Configurable variable environments in the chart values are used to do this.

Parameters            | Default values                            | Comment
---                   | ---                                       | --- 
Namespace             | "prescaling-exporter"                     | Namespace for PrescalingEvent
Port                  | "9101"                                    | Application port
AnnotationEndTime     | "annotations.scaling.exporter.time.end"   | Annotation end in HPA for create metrique
AnnotationStartTime   | "annotations.scaling.exporter.time.start" | Annotation start in HPA for create metrique
AnnotationMinReplicas | "annotations.scaling.exporter.replica.min"| Annotation in min HPA for create metrique
LabelProject          | "project"                                 | k8s label used to add a label in the prescaling metric

# Prescaling Web server
## OpenAPI docs

In the application is provided a swagger which is accessible with the endpoint. 

/swagger/index.html

## Event API 

To allow the platform to scale up and down on different schedules and with a multiplier, pre-scaling events can be requested by the DRBs or the API.

The API allows the creation/modification and deletion of Prescaling Event CRDs in the cluster

/api/v1/events/

## Metrics

The metrics are exposed with the endpoint

/metrics

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
