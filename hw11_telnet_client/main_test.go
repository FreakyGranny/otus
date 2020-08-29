package main

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TelnetSuite struct {
	suite.Suite
	mockTelnetCtl *gomock.Controller
	mockTelnet    *MockTelnetClient
}

func (s *TelnetSuite) SetupTest() {
	s.mockTelnetCtl = gomock.NewController(s.T())
	s.mockTelnet = NewMockTelnetClient(s.mockTelnetCtl)
}

func (s *TelnetSuite) TearDownTest() {
	s.mockTelnetCtl.Finish()
}

func (s *TelnetSuite) TestOk() {
	s.mockTelnet.EXPECT().Connect().Return(nil)
	s.mockTelnet.EXPECT().Send().Return(nil)
	s.mockTelnet.EXPECT().Receive().Return(nil)
	s.mockTelnet.EXPECT().Close().Return(nil)

	s.Require().NoError(run(s.mockTelnet))
}

func (s *TelnetSuite) TestConnectFail() {
	s.mockTelnet.EXPECT().Connect().Return(errors.New("fake error"))

	s.Require().Error(run(s.mockTelnet))
}

func TestTelnetSuite(t *testing.T) {
	suite.Run(t, new(TelnetSuite))
}
