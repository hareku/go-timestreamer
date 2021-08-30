package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hareku/go-timestreamer/simple"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

var recorder simple.Recorder
var nowFunc func() time.Time

func init() {
	nowFunc = time.Now
	recorder = simple.NewTimestreamRecorder(&simple.NewTimestreamRecorderInput{
		TimestreamWrite: timestreamwrite.New(session.Must(session.NewSession())),
		DatabaseName:    os.Getenv("DATABASE_NAME"),
		TableName:       os.Getenv("TABLE_NAME"),
		MeasureName:     os.Getenv("MEASURE_NAME"),
	})
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := recorder.Do(ctx, nowFunc())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("recording failed: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
