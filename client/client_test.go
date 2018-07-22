package client

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	opts Options
}

func (s *ClientTestSuite) SetupTest() {
	domain := "server.com"
	token := "a_token"
	s.opts = Options{Domain: &domain, Token: &token}
}

func (s *ClientTestSuite) TearDownTest() {
	gock.Off()

	os.Unsetenv("DATABRICKS_DOMAIN")
	os.Unsetenv("DATABRICKS_TOKEN")
}

func (s *ClientTestSuite) TestNewClientFailsIfMissingCredentials() {
	_, err := NewClient(Options{})
	s.Assert().EqualError(err, "missing credentials")
}

func (s *ClientTestSuite) TestNewClientLoadsCredentialsFromEnvironment() {
	os.Setenv("DATABRICKS_DOMAIN", "server.com")
	os.Setenv("DATABRICKS_TOKEN", "a_token")

	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(200).
		BodyString("a response")

	cl, err := NewClient(Options{})
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestQueryHandlesResponseErrors() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(500).
		JSON(map[string]string{
			"error_code": "an error code",
			"message":    "a message",
		})

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().Error(err)

	s.Assert().IsType(Error{}, err)
	s.Assert().Equal(
		"an error code", err.(Error).Code())
	s.Assert().Equal(
		"a message", err.(Error).Error())
}

func (s *ClientTestSuite) TestQueryHandlesNonJsonErrors() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(500).
		BodyString("an error")

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().Error(err)

	s.Assert().Equal("request error: an error", err.Error())
}

func (s *ClientTestSuite) TestQueryReturnsSuccessfulResponse() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(200).
		BodyString("a response")

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	resp, err := cl.Query("GET", "foo", nil)
	s.Require().NoError(err)

	s.Assert().Equal("a response", string(resp))
}

func (s *ClientTestSuite) TestAuthHeaderIsPresent() {
	gock.New("https://server.com").
		MatchHeader("Authorization", "Bearer a_token").
		Get("^/api/2.0/foo$").
		Reply(200)

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	cl.Query("GET", "foo", nil)

	s.Assert().True(gock.IsDone())
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
