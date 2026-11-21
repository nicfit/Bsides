package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MainTests struct {
	suite.Suite
	stdout, stderr *bytes.Buffer
	run            func(args []string) error
}

func (s *MainTests) SetupSuite() {}

func (s *MainTests) TearDownSuite() {}

func (s *MainTests) SetupTest() {
	s.stdout = new(bytes.Buffer)
	s.stderr = new(bytes.Buffer)
	s.run = func(args []string) error {
		return Main(args, s.stdout, s.stderr)
	}
}

func (s *MainTests) TearDownTest() {}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(MainTests))
}

func (s *MainTests) TestNoArgs() {
	err := s.run([]string{})
	assert.Error(s.T(), err)
	assert.ErrorContains(s.T(), err, "usage:")
}

func (s *MainTests) TestFile() {
	err := s.run([]string{"main_test.go"})
	assert.NoError(s.T(), err)
}
