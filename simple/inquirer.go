package simple

import (
	"context"
	"time"
)

//go:generate mockgen -source=${GOFILE} -destination=./${GOPACKAGE}mock/mock_${GOFILE} -package=${GOPACKAGE}mock

type RecordOfDay struct {
	Date time.Time
	Sum  uint64
}

type PerDayInput struct {
	Days         uint32
	Location     *time.Location
	DatabaseName string
	TableName    string
}

type Inquirer interface {
	PerDay(ctx context.Context, in PerDayInput) ([]RecordOfDay, error)
}
