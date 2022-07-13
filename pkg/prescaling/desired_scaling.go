package prescaling

const (
	UPSCALE     int = 11
	NODOWNSCALE     = 10
	DOWNSCALE       = 9
	NOSCALE         = 0
)

func DesiredScaling(eventInRangeTime bool, multiplier int, minReplica int, currentReplica int32) int {
	if !eventInRangeTime {
		return NOSCALE
	}

	if multiplier == 0 {
		multiplier = 1
	}

	desiredReplica := multiplier * minReplica

	if desiredReplica > int(currentReplica) {
		return UPSCALE
	}

	if desiredReplica == int(currentReplica) {
		return NODOWNSCALE
	}

	threshold := ((float64(currentReplica) - float64(desiredReplica)) / float64(desiredReplica)) * float64(100)
	if threshold > float64(10) {
		return DOWNSCALE
	}

	return NODOWNSCALE
}
