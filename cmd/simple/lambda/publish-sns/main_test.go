package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/hareku/go-timestreamer/internal/snssdk/snssdkmock"
	"github.com/hareku/go-timestreamer/simple"
	"github.com/hareku/go-timestreamer/simple/simplemock"
	"github.com/stretchr/testify/require"
)

func Test_handler_Successful(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	inqMock := simplemock.NewMockInquirer(ctrl)
	inqMock.EXPECT().PerDay(ctx, simple.PerDayInput{
		Days:         30,
		Location:     loc,
		DatabaseName: "mydatabase",
		TableName:    "mytable",
	}).Times(1).Return([]simple.RecordOfDay{
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
	}, nil)

	snsMock := snssdkmock.NewMockSNSIface(ctrl)
	snsMock.EXPECT().PublishWithContext(ctx, gomock.Eq(&sns.PublishInput{
		TopicArn: aws.String("arn12345"),
		Subject:  aws.String("Simple reporting by go-timestreamer"),
		Message:  aws.String("2020-08-15: 100\n2020-08-14: 80\n"),
	})).Times(1).Return(nil, nil)

	// set global vars
	inquirer = inqMock
	snsapi = snsMock
	databaseName = "mydatabase"
	tableName = "mytable"
	snsTopicArn = "arn12345"
	location = loc

	err = handler(ctx)
	require.NoError(t, err)
}

func Test_handler_SNSFailed(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	inqMock := simplemock.NewMockInquirer(ctrl)
	inqMock.EXPECT().PerDay(ctx, simple.PerDayInput{
		Days:         30,
		Location:     loc,
		DatabaseName: "mydatabase",
		TableName:    "mytable",
	}).Times(1).Return([]simple.RecordOfDay{
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
	}, nil)

	snsMock := snssdkmock.NewMockSNSIface(ctrl)
	snsMock.EXPECT().PublishWithContext(ctx, gomock.Eq(&sns.PublishInput{
		TopicArn: aws.String("arn12345"),
		Subject:  aws.String("Simple reporting by go-timestreamer"),
		Message:  aws.String("2020-08-15: 100\n2020-08-14: 80\n"),
	})).Times(1).Return(nil, errors.New("sns error"))

	// set global vars
	inquirer = inqMock
	snsapi = snsMock
	databaseName = "mydatabase"
	tableName = "mytable"
	snsTopicArn = "arn12345"
	location = loc

	err = handler(ctx)
	require.Error(t, err)
}

func Test_handler_InquiringFailed(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	inqMock := simplemock.NewMockInquirer(ctrl)
	inqMock.EXPECT().PerDay(ctx, simple.PerDayInput{
		Days:         30,
		Location:     loc,
		DatabaseName: "mydatabase",
		TableName:    "mytable",
	}).Times(1).Return(nil, errors.New("timestream error"))

	// set global vars
	inquirer = inqMock
	databaseName = "mydatabase"
	tableName = "mytable"
	snsTopicArn = "arn12345"
	location = loc

	err = handler(ctx)
	require.Error(t, err)
}
