package simple

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/hareku/go-timestreamer/internal/timestreamsdk"
)

type NewTimestreamRecorderInput struct {
	TimestreamWrite timestreamsdk.TimestreamWriteIface
	DatabaseName    string
	TableName       string
	MeasureName     string
}

func NewTimestreamRecorder(in *NewTimestreamRecorderInput) Recorder {
	return &timstreamRecorder{
		timestreamWrite: in.TimestreamWrite,
		databaseName:    in.DatabaseName,
		tableName:       in.TableName,
		measureName:     in.MeasureName,
	}
}

type timstreamRecorder struct {
	timestreamWrite timestreamsdk.TimestreamWriteIface
	databaseName    string
	tableName       string
	measureName     string
}

func (r *timstreamRecorder) Do(ctx context.Context, t time.Time) error {
	_, err := r.timestreamWrite.WriteRecordsWithContext(ctx, &timestreamwrite.WriteRecordsInput{
		DatabaseName: &r.databaseName,
		TableName:    &r.tableName,
		Records: []*timestreamwrite.Record{
			{
				MeasureName:      &r.measureName,
				MeasureValue:     aws.String("1"),
				MeasureValueType: aws.String("BIGINT"),
				Time:             aws.String(fmt.Sprintf("%d", t.UnixNano())),
				TimeUnit:         aws.String("NANOSECONDS"),
				Dimensions: []*timestreamwrite.Dimension{
					{
						Name:  aws.String("N"),
						Value: aws.String("N"),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("timestream failed to write records: %w", err)
	}
	return nil
}
