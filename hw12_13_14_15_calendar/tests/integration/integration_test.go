// +build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	internalgrpc "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type IntegrationSuite struct {
	suite.Suite
	cli internalgrpc.EventsClient
}

func (s *IntegrationSuite) Init(apiURL string) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(apiURL, opts...)
	if err != nil {
		s.T().Fail()
	}

	s.cli = internalgrpc.NewEventsClient(conn)
}

func (s *IntegrationSuite) SetupTest() {
	apiURL := os.Getenv("GRPC_URL")
	if apiURL == "" {
		apiURL = "127.0.0.1:50051"
	}

	s.Init(apiURL)
}

func (s *IntegrationSuite) TestValidCreate() {
	startDate, _ := ptypes.TimestampProto(time.Date(2020, 9, 25, 10, 41, 0, 0, time.UTC))
	endDate, _ := ptypes.TimestampProto(time.Date(2020, 9, 25, 13, 0, 0, 0, time.UTC))
	request := &internalgrpc.CreateEventRequest{
		Title:        "test",
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        "some descr",
		OwnerId:      1,
		NotifyBefore: 0,
	}
	response, err := s.cli.CreateEvent(context.Background(), request)

	if err != nil {
		s.T().Fail()
	}
	s.Require().NotNil(response)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
