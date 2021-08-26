package timestreamsdk

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

//go:generate mockgen -source=${GOFILE} -destination=./${GOPACKAGE}mock/mock_${GOFILE} -package=${GOPACKAGE}mock

type TimestreamWriteIface interface {
	WriteRecordsWithContext(ctx context.Context, input *timestreamwrite.WriteRecordsInput, opts ...request.Option) (*timestreamwrite.WriteRecordsOutput, error)
}
