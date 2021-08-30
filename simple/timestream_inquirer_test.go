package simple_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/golang/mock/gomock"
	"github.com/hareku/go-timestreamer/internal/timestreamsdk/timestreamsdkmock"
	"github.com/hareku/go-timestreamer/simple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_timestreamInquirer_PerDay(t *testing.T) {
	t.Parallel()

	t.Run("Successful", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		sdk := timestreamsdkmock.NewMockTimestreamQueryIface(ctrl)

		sdk.EXPECT().QueryPagesWithContext(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, in interface{}, fn func(page *timestreamquery.QueryOutput, lastPage bool) bool, opts ...interface{}) error {
				fn(&timestreamquery.QueryOutput{
					Rows: []*timestreamquery.Row{
						{
							Data: []*timestreamquery.Datum{
								{
									ScalarValue: aws.String("2020-08-15 00:00:00.000000000"),
								},
								{
									ScalarValue: aws.String("100"),
								},
							},
						},
						{
							Data: []*timestreamquery.Datum{
								{
									ScalarValue: aws.String("2020-08-14 00:00:00.000000000"),
								},
								{
									ScalarValue: aws.String("80"),
								},
							},
						},
					},
				}, true)
				return nil
			})

		loc, err := time.LoadLocation("Asia/Tokyo")
		require.NoError(t, err)

		inq := simple.NewTimestreamInquire(sdk)
		got, err := inq.PerDay(ctx, simple.PerDayInput{
			Location: loc,
		})
		require.NoError(t, err)
		assert.Equal(t, []simple.RecordOfDay{
			{
				Date: func() time.Time {
					res, err := time.ParseInLocation("2006-01-02", "2020-08-15", loc)
					require.NoError(t, err)
					return res
				}(),
				Sum: 100,
			},
			{
				Date: func() time.Time {
					res, err := time.ParseInLocation("2006-01-02", "2020-08-14", loc)
					require.NoError(t, err)
					return res
				}(),
				Sum: 80,
			},
		}, got)
	})

	t.Run("Timestream error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		sdk := timestreamsdkmock.NewMockTimestreamQueryIface(ctrl)

		sdk.EXPECT().QueryPagesWithContext(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("unexpected error"))

		loc, err := time.LoadLocation("Asia/Tokyo")
		require.NoError(t, err)

		inq := simple.NewTimestreamInquire(sdk)
		_, err = inq.PerDay(ctx, simple.PerDayInput{
			Location: loc,
		})
		require.Error(t, err)
	})
}
