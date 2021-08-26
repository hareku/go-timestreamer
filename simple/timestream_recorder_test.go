package simple

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hareku/go-timestreamer/internal/timestreamsdk/timestreamsdkmock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_timstreamRecorder_Do(t *testing.T) {
	t.Parallel()

	t.Run("Successful write Timestream", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		tw := timestreamsdkmock.NewMockTimestreamWriteIface(ctrl)
		tw.EXPECT().
			WriteRecordsWithContext(gomock.Eq(ctx), gomock.Eq(&timestreamwrite.WriteRecordsInput{
				DatabaseName: aws.String("database"),
				TableName:    aws.String("table"),
				Records: []*timestreamwrite.Record{
					{
						MeasureName:      aws.String("recorder"),
						MeasureValue:     aws.String("1"),
						MeasureValueType: aws.String("BIGINT"),
						Time:             aws.String("12345"),
						TimeUnit:         aws.String("NANOSECONDS"),
						Dimensions: []*timestreamwrite.Dimension{
							{
								Name:  aws.String("N"),
								Value: aws.String("N"),
							},
						},
					},
				},
			})).
			Times(1).
			Return(nil, nil)

		r := NewTimestreamRecorder(&NewTimestreamRecorderInput{
			TimestreamWrite: tw,
			DatabaseName:    "database",
			TableName:       "table",
			MeasureName:     "recorder",
		})
		err := r.Do(ctx, time.Unix(0, 12345))
		require.NoError(t, err)
	})

	t.Run("Unable to write Timestream", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		tw := timestreamsdkmock.NewMockTimestreamWriteIface(ctrl)
		tw.EXPECT().
			WriteRecordsWithContext(gomock.Eq(ctx), gomock.Eq(&timestreamwrite.WriteRecordsInput{
				DatabaseName: aws.String("database"),
				TableName:    aws.String("table"),
				Records: []*timestreamwrite.Record{
					{
						MeasureName:      aws.String("recorder"),
						MeasureValue:     aws.String("1"),
						MeasureValueType: aws.String("BIGINT"),
						Time:             aws.String("12345"),
						TimeUnit:         aws.String("NANOSECONDS"),
						Dimensions: []*timestreamwrite.Dimension{
							{
								Name:  aws.String("N"),
								Value: aws.String("N"),
							},
						},
					},
				},
			})).
			Times(1).
			Return(nil, errors.New("something happend"))

		r := NewTimestreamRecorder(&NewTimestreamRecorderInput{
			TimestreamWrite: tw,
			DatabaseName:    "database",
			TableName:       "table",
			MeasureName:     "recorder",
		})
		err := r.Do(ctx, time.Unix(0, 12345))
		require.Error(t, err)
	})
}
