package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/hareku/go-timestreamer/internal/snssdk"
	"github.com/hareku/go-timestreamer/simple"
)

var inquirer simple.Inquirer
var snsapi snssdk.SNSIface
var databaseName, tableName, snsTopicArn string
var location *time.Location

func init() {
	sess := session.Must(session.NewSession())

	inquirer = simple.NewTimestreamInquire(timestreamquery.New(sess))
	snsapi = sns.New(sess)

	databaseName = os.Getenv("DATABASE_NAME")
	tableName = os.Getenv("TABLE_NAME")
	snsTopicArn = os.Getenv("SNS_TOPIC_ARN")

	tz := os.Getenv("GO_TIMEZONE")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(fmt.Errorf("loading location %q failed: %w", tz, err))
	}
	location = loc
}

func handler(ctx context.Context) error {
	records, err := inquirer.PerDay(ctx, simple.PerDayInput{
		Days:         30,
		Location:     location,
		DatabaseName: databaseName,
		TableName:    tableName,
	})
	if err != nil {
		return err
	}

	msg := ""
	for _, rec := range records {
		msg += fmt.Sprintf("%s: %d\n", rec.Date.Format("2006-01-02"), rec.Sum)
	}

	_, err = snsapi.PublishWithContext(ctx, &sns.PublishInput{
		TopicArn: aws.String(snsTopicArn),
		Subject:  aws.String("Simple reporting by go-timestreamer"),
		Message:  aws.String(msg),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
