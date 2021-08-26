package simple

import (
	"context"
	"time"
)

//go:generate mockgen -source=${GOFILE} -destination=./${GOPACKAGE}mock/mock_${GOFILE} -package=${GOPACKAGE}mock

// Recorder provides simple recording functions for Amazon Timestream.
type Recorder interface {
	// Do writes one time series data.
	// MeasureValue: 1
	// MeasureValueType: BIGINT
	// TimeUnit: NANOSECONDS
	// Dimensions[0].Name: N
	// Dimensions[0].Value: N
	Do(ctx context.Context, t time.Time) error
}
