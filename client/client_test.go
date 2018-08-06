package client

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"net"
	"net/url"
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
	s.opts = Options{Domain: &domain, Token: &token, MaxRetries: 0, RetryDelay: 0}
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

func (s *ClientTestSuite) TestQueryRetriesNetworkErrors() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Times(3).
		ReplyError(errorMock{temporary: true})

	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(200).
		BodyString("a response")

	s.opts.MaxRetries = 3

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestQueryRetriesDatabrickErrors() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Times(3).
		Reply(500).
		JSON(map[string]string{
			"error_code": "an error code",
			"message":    "a message",
		})

	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Reply(200).
		BodyString("a response")

	s.opts.MaxRetries = 3

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().NoError(err)
}

func (s *ClientTestSuite) TestQueryDoesNotRetryPermanentErrors() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		ReplyError(&net.DNSError{})

	s.opts.MaxRetries = 1

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().Error(err)

	urlErr, ok := err.(*url.Error)
	s.Require().True(ok)

	s.Assert().IsType(&net.DNSError{}, urlErr.Err)
}

func (s *ClientTestSuite) TestQueryDoesNotRetryClientErrorsOrRedirects() {
	errorCodes := []int{301, 302, 400, 401, 403, 404}
	for _, code := range errorCodes {
		resp := gock.New("https://server.com").
			Get("^/api/2.0/foo$").
			Reply(code).
			BodyString("body")

		if code >= 300 && code < 400 {
			resp.AddHeader("Location", "none")
		}

		s.opts.MaxRetries = 1

		cl, err := NewClient(s.opts)
		s.Require().NoError(err)

		_, err = cl.Query("GET", "foo", nil)
		s.Require().Error(err)

		databricksErr, ok := err.(Error)
		s.Require().True(ok)

		s.Assert().Equal(code, databricksErr.statusCode)
	}
}

func (s *ClientTestSuite) TestQueryPropagatesErrorAfterMaxRetries() {
	gock.New("https://server.com").
		Get("^/api/2.0/foo$").
		Times(4).
		ReplyError(errorMock{temporary: true})

	s.opts.MaxRetries = 3

	cl, err := NewClient(s.opts)
	s.Require().NoError(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Require().Error(err)

	s.Assert().Contains(err.Error(), "errorMock")
}

type errorMock struct {
	temporary bool
}

func (e errorMock) Error() string {
	return "errorMock"
}

func (e errorMock) Temporary() bool {
	return e.temporary
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
