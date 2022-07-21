package config

import (
	"os"
)

type Struct struct {
	Namespace             string
	Port                  string
	AnnotationEndTime     string
	AnnotationStartTime   string
	AnnotationMinReplicas string
	LabelProject          string
}

var Config = Struct{
	Namespace:             "prescaling-exporter",
	Port:                  "9101",
	AnnotationEndTime:     "annotations.scaling.exporter.time.end",
	AnnotationStartTime:   "annotations.scaling.exporter.time.start",
	AnnotationMinReplicas: "annotations.scaling.exporter.replica.min",
	LabelProject:          "project",
}

func init() {
	if os.Getenv("NAMESPACE") != "" {
		Config.Namespace = os.Getenv("NAMESPACE")
	}
	if os.Getenv("PORT") != "" {
		Config.Port = os.Getenv("PORT")
	}
	if os.Getenv("ANNOTATION_END_TIME") != "" {
		Config.AnnotationEndTime = os.Getenv("ANNOTATION_END_TIME")
	}
	if os.Getenv("ANNOTATION_START_TIME") != "" {
		Config.AnnotationStartTime = os.Getenv("ANNOTATION_START_TIME")
	}
	if os.Getenv("ANNOTATION_MIN_REPLICAS") != "" {
		Config.AnnotationMinReplicas = os.Getenv("ANNOTATION_MIN_REPLICAS")
	}
	if os.Getenv("LABEL_PROJECT") != "" {
		Config.LabelProject = os.Getenv("LABEL_PROJECT")
	}
}
