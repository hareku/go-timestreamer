package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hareku/go-timestreamer/simple/simplemock"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
)

func TestHandler_StatusOK(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	nowFunc = func() time.Time {
		return now
	}

	ctrl := gomock.NewController(t)
	mock := simplemock.NewMockRecorder(ctrl)
	mock.EXPECT().Do(ctx, now).Times(1).Return(nil)
	recorder = mock

	resp, err := handler(ctx, events.APIGatewayProxyRequest{})
	if err != nil {
		t.Errorf("Got error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTP Status Code is not 200, got %d", resp.StatusCode)
	}
}

func TestHandler_StatusInternalServerError(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	nowFunc = func() time.Time {
		return now
	}

	ctrl := gomock.NewController(t)
	mock := simplemock.NewMockRecorder(ctrl)
	mock.EXPECT().Do(ctx, now).Times(1).Return(errors.New("something happend"))
	recorder = mock

	resp, err := handler(ctx, events.APIGatewayProxyRequest{})
	if err == nil {
		t.Error("Got nil error")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("HTTP Status Code is not 500, got %d", resp.StatusCode)
	}
}
