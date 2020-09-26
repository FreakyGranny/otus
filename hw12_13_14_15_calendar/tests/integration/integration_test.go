// +build integration

package integration_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	internalgrpc "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/bxcodec/faker/v3"
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

func (s *IntegrationSuite) CreateEvent(date time.Time) int64 {
	s.T().Helper()
	startDate, _ := ptypes.TimestampProto(date)
	endDate, _ := ptypes.TimestampProto(date.Add(time.Hour))
	request := &internalgrpc.CreateEventRequest{
		Title:        faker.Username(),
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        faker.Sentence(),
		OwnerId:      rand.Int63n(10000000) + 1,
		NotifyBefore: rand.Int63n(1000),
	}
	response, err := s.cli.CreateEvent(context.Background(), request)
	s.Require().Nil(err)
	s.Require().NotNil(response)
	s.Require().NotEqual(0, response.Id)

	return response.Id
}

func (s *IntegrationSuite) DropEvent(id int64) {
	s.T().Helper()
	cleanupRequest := &internalgrpc.DeleteEventRequest{
		Id: id,
	}
	_, err := s.cli.DeleteEvent(context.Background(), cleanupRequest)
	s.Require().Nil(err)
}

func buildListRequest(date time.Time) *internalgrpc.ListEventRequest {
	reqDate, _ := ptypes.TimestampProto(date)

	return &internalgrpc.ListEventRequest{Date: reqDate}
}

func getMiddleDay() time.Time {
	var date time.Time
	now := time.Now()
	for i := 15; ; i++ {
		date = time.Date(now.Year(), now.Month(), i, 0, 0, 0, 0, time.UTC)
		if date.Weekday() > 1 && date.Weekday() < 6 {
			break
		}
	}

	return date
}
func (s *IntegrationSuite) TestValidCreate() {
	now := time.Now().UTC()
	nowShifted := now.Add(time.Hour)
	startDate, _ := ptypes.TimestampProto(now)
	endDate, _ := ptypes.TimestampProto(nowShifted)
	request := &internalgrpc.CreateEventRequest{
		Title:        "test",
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        "some descr",
		OwnerId:      111,
		NotifyBefore: 123456,
	}
	response, err := s.cli.CreateEvent(context.Background(), request)
	s.Require().Nil(err)
	s.Require().NotNil(response)
	s.Require().NotEqual(0, response.GetId())
	s.Require().Equal("test", response.GetTitle())
	s.Require().Equal(now, response.GetStartDate().AsTime())
	s.Require().Equal(nowShifted, response.GetEndDate().AsTime())
	s.Require().Equal("some descr", response.GetDescr())
	s.Require().Equal(int64(111), response.GetOwnerId())
	s.Require().Equal(int64(123456), response.GetNotifyBefore())

	s.DropEvent(response.Id)
}

func (s *IntegrationSuite) TestCreateInvalidTitle() {
	startDate, _ := ptypes.TimestampProto(time.Now())
	endDate, _ := ptypes.TimestampProto(time.Now().Add(time.Hour))
	request := &internalgrpc.CreateEventRequest{
		Title:        "",
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        "some descr",
		OwnerId:      1,
		NotifyBefore: 0,
	}
	_, err := s.cli.CreateEvent(context.Background(), request)
	s.Require().NotNil(err)
}

func (s *IntegrationSuite) TestCreateInvalidDates() {
	startDate, _ := ptypes.TimestampProto(time.Now())
	endDate, _ := ptypes.TimestampProto(time.Now().Add(-1 * time.Second))
	request := &internalgrpc.CreateEventRequest{
		Title:        "test",
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        "some descr",
		OwnerId:      1,
		NotifyBefore: 0,
	}
	_, err := s.cli.CreateEvent(context.Background(), request)
	s.Require().NotNil(err)
}

func (s *IntegrationSuite) TestListEmpty() {
	reqDate, _ := ptypes.TimestampProto(time.Now())
	listRequest := &internalgrpc.ListEventRequest{
		Date: reqDate,
	}
	response, err := s.cli.GetEventForDay(context.Background(), listRequest)
	s.Require().Nil(err)
	s.Require().Equal(0, len(response.Results))
}

func (s *IntegrationSuite) TestListInvalidWeekday() {
	_, err := s.cli.GetEventForWeek(context.Background(), buildListRequest(getMiddleDay()))
	s.Require().NotNil(err)
}

func (s *IntegrationSuite) TestListInvalidMonthDay() {
	_, err := s.cli.GetEventForMonth(context.Background(), buildListRequest(getMiddleDay()))
	s.Require().NotNil(err)
}

func (s *IntegrationSuite) TestList() {
	testingDate := getMiddleDay()
	dateAtDay := time.Date(testingDate.Year(), testingDate.Month(), testingDate.Day(), 0, 0, 0, 0, time.UTC)
	monday := testingDate.Day() - int(testingDate.Weekday()) + 1
	dateAtWeek := time.Date(testingDate.Year(), testingDate.Month(), monday, 0, 0, 0, 0, time.UTC)
	dateAtMonth := time.Date(testingDate.Year(), testingDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	id1 := s.CreateEvent(dateAtDay.Add(time.Minute))      // current day
	id2 := s.CreateEvent(dateAtDay.Add(time.Hour * 50))   // on current week
	id3 := s.CreateEvent(dateAtDay.Add(time.Hour * 200))  // on current month
	id4 := s.CreateEvent(dateAtMonth.Add(time.Hour * -1)) // on previous month

	response, err := s.cli.GetEventForDay(context.Background(), buildListRequest(dateAtDay))
	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Results))

	response, err = s.cli.GetEventForWeek(context.Background(), buildListRequest(dateAtWeek))
	s.Require().Nil(err)
	s.Require().Equal(2, len(response.Results))

	response, err = s.cli.GetEventForMonth(context.Background(), buildListRequest(dateAtMonth))
	s.Require().Nil(err)
	s.Require().Equal(3, len(response.Results))

	s.DropEvent(id1)
	s.DropEvent(id2)
	s.DropEvent(id3)
	s.DropEvent(id4)
}

func (s *IntegrationSuite) TestUpdate() {
	eventId := s.CreateEvent(time.Now())
	startDate, _ := ptypes.TimestampProto(time.Now())
	endDate, _ := ptypes.TimestampProto(time.Now().Add(time.Hour))

	request := &internalgrpc.UpdateEventRequest{
		Id:           eventId,
		Title:        "test",
		StartDate:    startDate,
		EndDate:      endDate,
		Descr:        "",
		OwnerId:      0,
		NotifyBefore: 0,
	}
	response, err := s.cli.UpdateEvent(context.Background(), request)
	s.Require().Nil(err)
	s.Require().NotNil(response)
	s.Require().Equal(eventId, response.GetId())
	s.Require().Equal("test", response.GetTitle())

	s.DropEvent(response.Id)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
