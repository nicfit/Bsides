package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/nicfit/bsides/factory"
	"github.com/nicfit/bsides/m3u"
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
		fmt.Println("Running _main:", args)
		args = append([]string{"_main"}, args...)
		return _main(args, s.stdout, s.stderr)
	}
}

func (s *MainTests) TearDownTest() {}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(MainTests))
}

func (s *MainTests) TestNoArgs() {
	err := s.run([]string{})
	assert.Error(s.T(), err)
	assert.ErrorContains(s.T(), err, "Required argument \"playlist\" not set")
}

func (s *MainTests) TestFile() {
	err := s.run([]string{"main_test.go"})
	assert.Error(s.T(), err)
}

func (s *MainTests) TestFactoryLoad() {
	tmpfiles := make([]string, 0)
	defer func() {
		for i, _ := range tmpfiles {
			_ = os.Remove(tmpfiles[i])
		}
	}()

	type testcase struct {
		ext            string
		expect_success bool
		invalid        bool
	}

	for _, tc := range []testcase{
		{ext: m3u.DOT_M3U, expect_success: true},
		{ext: m3u.DOT_M3U8, expect_success: true},
		{ext: ".pls", expect_success: false},
		{ext: ".xspf", expect_success: false},
		{ext: ".xml", expect_success: false},
		{ext: ".jspf", expect_success: false},
		{ext: ".json", expect_success: false},
		{ext: ".foo", expect_success: false, invalid: true},
	} {
		tmp, err := os.CreateTemp("", "pl*"+tc.ext)
		assert.NoError(s.T(), err)
		tmpfiles = append(tmpfiles, tmp.Name())
		playlist, err := factory.Load(tmp.Name())
		if tc.expect_success {
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), playlist)
		} else {
			assert.Error(s.T(), err)
			if tc.invalid {
				assert.ErrorContains(s.T(), err, "invalid file extension: "+tc.ext)
			} else {
				assert.ErrorContains(s.T(), err, "unsupported file extension: "+tc.ext)
			}
			assert.Nil(s.T(), playlist)
		}
	}
}

func (s *MainTests) TestDefault() {
	assert.NoError(s.T(), s.run([]string{"../m3u/sample.m3u"}))
	scanner := bufio.NewScanner(s.stdout)
	for i := range 5 {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			assert.Equal(s.T(), fmt.Sprintf("file%d.mp3", i+1), line)
		}
	}
	assert.NoError(s.T(), scanner.Err())
}

func (s *MainTests) TestIter() {
	err := s.run([]string{"iter", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)
	scanner := bufio.NewScanner(s.stdout)
	for i := range 5 {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		assert.NoError(s.T(), scanner.Err())
		if line != "" {
			assert.Equal(s.T(), fmt.Sprintf("file%d.mp3", i+1), line)
		}
	}
}

func (s *MainTests) TestIterReverse() {
	err := s.run([]string{"iter", "--reverse", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)
	scanner := bufio.NewScanner(s.stdout)
	non_blank := false
	for i := 5; i > 0; i-- {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		assert.NoError(s.T(), scanner.Err())
		if line != "" {
			assert.Equal(s.T(), fmt.Sprintf("file%d.mp3", i), line)
			non_blank = true
		}
	}
	assert.True(s.T(), non_blank)
}

func (s *MainTests) TestIterRepeat() {
	err := s.run([]string{"iter", "--repeat=3", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)
	scanner := bufio.NewScanner(s.stdout)
	non_blank := false
	for _ = range 3 {
		for i := range 5 {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			assert.NoError(s.T(), scanner.Err())
			if line != "" {
				assert.Equal(s.T(), fmt.Sprintf("file%d.mp3", i+1), line)
				non_blank = true
			}
		}
	}
	assert.True(s.T(), non_blank)
}

func (s *MainTests) TestIterReverseRepeat() {
	err := s.run([]string{"iter", "--reverse", "--repeat=5", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)
	scanner := bufio.NewScanner(s.stdout)
	non_blank := false
	for _ = range 5 {
		for i := 5; i > 0; i-- {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			assert.NoError(s.T(), scanner.Err())
			if line != "" {
				assert.Equal(s.T(), fmt.Sprintf("file%d.mp3", i), line)
				non_blank = true
			}
		}
	}
	assert.True(s.T(), non_blank)
}

func (s *MainTests) TestShuffle() {
	err := s.run([]string{"iter", "--shuffle", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)

	chosen := make([]string, 0)
	scanner := bufio.NewScanner(s.stdout)
	for _ = range 5 {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		assert.NoError(s.T(), scanner.Err())
		if line != "" {
			assert.NotContains(s.T(), chosen, line)
			assert.True(s.T(), filenameRegex.MatchString(line))
			chosen = append(chosen, line)
		}
	}
}

func (s *MainTests) TestShuffleReverse() {
	err := s.run([]string{"iter", "--shuffle", "--reverse", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)

	chosen := make([]string, 0)
	scanner := bufio.NewScanner(s.stdout)
	for _ = range 4 {
		for _ = range 5 {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			assert.NoError(s.T(), scanner.Err())
			if line != "" {
				assert.NotContains(s.T(), chosen, line)
				assert.True(s.T(), filenameRegex.MatchString(line))
				chosen = append(chosen, line)
			}
		}
		chosen = make([]string, 0)
	}
}

func (s *MainTests) TestShuffleRepeat() {
	err := s.run([]string{"iter", "--shuffle", "--repeat=7", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)

	chosen := make([]string, 0)
	scanner := bufio.NewScanner(s.stdout)
	for _ = range 7 {
		for _ = range 5 {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			assert.NoError(s.T(), scanner.Err())
			if line != "" {
				assert.NotContains(s.T(), chosen, line)
				assert.True(s.T(), filenameRegex.MatchString(line))
				chosen = append(chosen, line)
			}
		}
		chosen = make([]string, 0)
	}
}

func (s *MainTests) TestShuffleRepeatReverse() {
	err := s.run([]string{"iter", "--shuffle", "--reverse", "--repeat=4", "../m3u/sample-extended.m3u"})
	assert.NoError(s.T(), err)

	chosen := make([]string, 0)
	scanner := bufio.NewScanner(s.stdout)
	for _ = range 4 {
		for _ = range 5 {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			assert.NoError(s.T(), scanner.Err())
			if line != "" {
				assert.NotContains(s.T(), chosen, line)
				assert.True(s.T(), filenameRegex.MatchString(line))
				chosen = append(chosen, line)
			}
		}
		chosen = make([]string, 0)
	}
}

var filenameRegex = regexp.MustCompile(`file[1-5].mp3`)
