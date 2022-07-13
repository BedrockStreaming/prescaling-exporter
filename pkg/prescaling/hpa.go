package prescaling

import (
	"errors"
	"time"
)

type Hpa struct {
	Replica         int
	CurrentReplicas int32
	Start           time.Time
	End             time.Time
	Project         string
	Namespace       string
	Deployment      string
}

func (p Hpa) Validate() error {
	if p.Replica == 0 {
		err := errors.New("annotation replica min is misconfigured")
		return err
	}
	if p.Start.IsZero() {
		err := errors.New("annotation time start is misconfigured")
		return err
	}
	if p.End.IsZero() {
		err := errors.New("annotation time start is misconfigured")
		return err
	}

	return nil
}
