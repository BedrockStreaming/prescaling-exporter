package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {

	assert.Equal(t, Config.Port, "9101")
	assert.Equal(t, Config.AnnotationEndTime, "annotations.scaling.exporter.time.end")
	assert.Equal(t, Config.AnnotationStartTime, "annotations.scaling.exporter.time.start")
	assert.Equal(t, Config.LabelProject, "project")
}
