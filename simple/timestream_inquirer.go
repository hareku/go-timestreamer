package simple

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/hareku/go-timestreamer/internal/timestreamsdk"
)

func NewTimestreamInquire(sdk timestreamsdk.TimestreamQueryIface) Inquirer {
	return &timestreamInquirer{sdk}
}

type timestreamInquirer struct {
	sdk timestreamsdk.TimestreamQueryIface
}

func (i *timestreamInquirer) PerDay(ctx context.Context, in PerDayInput) ([]RecordOfDay, error) {
	_, offset := time.Now().In(in.Location).Zone()

	query := fmt.Sprintf(
		`SELECT bin(time + %ds, 24h) as date,
SUM(measure_value::bigint) as count
FROM "%s"."%s"
WHERE time between ago(%dd) and now()
GROUP BY bin(time + %ds, 24h)
ORDER BY date DESC`,
		offset,
		in.DatabaseName,
		in.TableName,
		in.Days,
		offset,
	)

	records := make([]RecordOfDay, 0)
	var parseErr error
	err := i.sdk.QueryPagesWithContext(ctx, &timestreamquery.QueryInput{
		QueryString: aws.String(query),
	}, func(page *timestreamquery.QueryOutput, lastPage bool) bool {
		rows := page.Rows
		for _, row := range rows {
			timestamp := *row.Data[0].ScalarValue
			date, err := time.ParseInLocation("2006-01-02 15:04:05.999999999", timestamp, in.Location)
			if err != nil {
				parseErr = fmt.Errorf("parsing date (%q) failed: %w", timestamp, err)
				return false
			}

			sumStr := *row.Data[1].ScalarValue
			sum, err := strconv.ParseUint(sumStr, 10, 64)
			if err != nil {
				parseErr = fmt.Errorf("parsing sum (%q) failed: %w", sumStr, err)
				return false
			}
			records = append(records, RecordOfDay{
				Date: date,
				Sum:  sum,
			})
		}

		return true
	})

	if err != nil {
		return nil, fmt.Errorf("timestream query failed: %w", err)
	}
	if parseErr != nil {
		return nil, fmt.Errorf("timestream parsing failed: %w", parseErr)
	}
	return records, nil
}
