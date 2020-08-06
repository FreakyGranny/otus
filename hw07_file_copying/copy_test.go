package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	TmpDirName string = "otus_tmp"
)

type CopySuite struct {
	suite.Suite
	TmpDirPath  string
	DstFilePath string
}

func (s *CopySuite) SetupTest() {
	s.TmpDirPath = path.Join(os.TempDir(), TmpDirName)
	os.Mkdir(s.TmpDirPath, 0755)
	s.DstFilePath = path.Join(s.TmpDirPath, "xxx")
}

func (s *CopySuite) TearDownTest() {
	os.RemoveAll(s.TmpDirPath)
}

func (s *CopySuite) TestCopy() {
	err := Copy("./testdata/input.txt", s.DstFilePath, 0, 0)
	s.Require().NoError(err)

	info, err := os.Stat("./testdata/input.txt")
	if err != nil {
		s.T().Fail()
	}
	resultInfo, err := os.Stat(s.DstFilePath)
	if err != nil {
		s.T().Fail()
	}
	s.Require().Equal(info.Size(), resultInfo.Size())
}

func (s *CopySuite) TestCopyWithLimit() {
	err := Copy("./testdata/input.txt", s.DstFilePath, 0, 10)
	s.Require().NoError(err)

	resultInfo, err := os.Stat(s.DstFilePath)
	if err != nil {
		s.T().Fail()
	}
	s.Require().Equal(int64(10), resultInfo.Size())	
}

func (s *CopySuite) TestCopyWithLimitOffset() {
	err := Copy("./testdata/input.txt", s.DstFilePath, 6000, 1000)
	s.Require().NoError(err)

	resultInfo, err := os.Stat(s.DstFilePath)
	if err != nil {
		s.T().Fail()
	}
	s.Require().Equal(int64(617), resultInfo.Size())
}

func (s *CopySuite) TestNoFile() {
	err := Copy("./there_is_no_file", s.DstFilePath, 0, 0)
	s.Require().Error(err)
	s.Require().EqualError(err, "stat ./there_is_no_file: no such file or directory")
}

func (s *CopySuite) TestSpecialFile() {
	err := Copy("/dev/null", s.DstFilePath, 0, 0)
	s.Require().Error(err)
	s.Require().Equal(err, ErrUnsupportedFile)
}

func (s *CopySuite) TestOffsetExceeds() {
	err := Copy("./testdata/input.txt", s.DstFilePath, 10000, 0)
	s.Require().Error(err)
	s.Require().Equal(err, ErrOffsetExceedsFileSize)
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(CopySuite))
}
