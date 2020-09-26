package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/mocks"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

type EventHandlerSuite struct {
	suite.Suite
	mockStorageCtl *gomock.Controller
	mockStorage    *mocks.MockStorage
}

func (s *EventHandlerSuite) SetupTest() {
	s.mockStorageCtl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockStorage(s.mockStorageCtl)
}

func (s *EventHandlerSuite) TearDownTest() {
	s.mockStorageCtl.Finish()
}

func (s *EventHandlerSuite) TestGetEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/event/1", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	e := &storage.Event{
		ID:    1,
		Title: "test event",
	}
	s.mockStorage.EXPECT().GetEvent(ctx, int64(1)).Return(e, nil)

	rec := httptest.NewRecorder()
	h := NewEventHandler(app.New(s.mockStorage))
	h.GetEvent(ctx, rec, req)
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `{"id":1,"title":"test event","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","descr":"","owner_id":0,"notify_before":0}`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *EventHandlerSuite) TestGetEventList() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/event?date=2020-09-26&period=d", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")
	e := []*storage.Event{
		{
			ID:    1,
			Title: "test event",
		},
		{
			ID:    2,
			Title: "second event",
		},
	}
	date := time.Date(2020, 9, 26, 0, 0, 0, 0, time.UTC)

	s.mockStorage.EXPECT().GetEventList(ctx, date, time.Hour*24).Return(e, nil)

	rec := httptest.NewRecorder()
	h := NewEventHandler(app.New(s.mockStorage))
	h.GetEventList(ctx, rec, req)
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `[{"id":1,"title":"test event","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","descr":"","owner_id":0,"notify_before":0},{"id":2,"title":"second event","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","descr":"","owner_id":0,"notify_before":0}]`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *EventHandlerSuite) TestCreateEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqEvent := &storage.Event{
		Title:        "new event",
		StartDate:    time.Date(2020, 9, 5, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2020, 9, 6, 0, 0, 0, 0, time.UTC),
		OwnerID:      1,
		Descr:        "description",
		NotifyBefore: 1,
	}
	body, err := json.Marshal(reqEvent)
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	s.mockStorage.EXPECT().CreateEvent(ctx, reqEvent).Return(nil)

	h := NewEventHandler(app.New(s.mockStorage))
	rec := httptest.NewRecorder()
	h.CreateEvent(ctx, rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code)

	var pJSON = `{"id":0,"title":"new event","start_date":"2020-09-05T00:00:00Z","end_date":"2020-09-06T00:00:00Z","descr":"description","owner_id":1,"notify_before":1}`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *EventHandlerSuite) TestUpdateEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqEvent := &storage.Event{
		Title: "new title",
	}
	body, err := json.Marshal(reqEvent)
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(http.MethodPatch, "/event/12", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	expect := &storage.Event{
		ID:           12,
		Title:        reqEvent.Title,
		StartDate:    reqEvent.StartDate,
		EndDate:      reqEvent.EndDate,
		OwnerID:      reqEvent.OwnerID,
		Descr:        reqEvent.Descr,
		NotifyBefore: reqEvent.NotifyBefore,
	}
	s.mockStorage.EXPECT().UpdateEvent(ctx, expect).Return(nil)

	h := NewEventHandler(app.New(s.mockStorage))
	rec := httptest.NewRecorder()
	h.UpdateEvent(ctx, rec, req)
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `{"id":12,"title":"new title","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","descr":"","owner_id":0,"notify_before":0}`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *EventHandlerSuite) TestDeleteProject() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := httptest.NewRequest(http.MethodDelete, "/event/18", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	s.mockStorage.EXPECT().DeleteEvent(ctx, int64(18)).Return(nil)

	rec := httptest.NewRecorder()
	h := NewEventHandler(app.New(s.mockStorage))
	h.DeleteEvent(ctx, rec, req)
	s.Require().Equal(http.StatusNoContent, rec.Code)
}

func TestEventHandlerSuite(t *testing.T) {
	suite.Run(t, new(EventHandlerSuite))
}
