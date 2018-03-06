package client

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
}

func (s *ClientTestSuite) TearDownTest() {
	gock.Off()
}

func (s *ClientTestSuite) TestQueryHandlesResponseErrors() {
	gock.New("https://server.com").
		Get("/api/2.0/foo").
		Reply(500).
		BodyString("a response").
		JSON(map[string]string{
			"error_code": "an error code",
			"message":    "a message",
		})

	cl, err := NewClient("server.com", "a_token")
	s.Assert().Nil(err)

	_, err = cl.Query("GET", "foo", nil)
	s.Assert().IsType(Error{}, err)
	s.Assert().Equal(
		"an error code", err.(Error).Code())
	s.Assert().Equal(
		"a message", err.(Error).Error())
}

func (s *ClientTestSuite) TestQueryReturnsSuccessfulResponse() {
	gock.New("https://server.com").
		Get("/api/2.0/foo").
		Reply(200).
		BodyString("a response")

	cl, err := NewClient("server.com", "a_token")
	s.Assert().Nil(err)

	resp, err := cl.Query("GET", "foo", nil)
	s.Assert().Nil(err)

	s.Assert().Equal("a response", string(resp))
}

func (s *ClientTestSuite) TestAuthHeaderIsPresent() {
	gock.New("https://server.com").
		MatchHeader("Authorization", "Bearer a_token").
		Get("/api/2.0/foo").
		Reply(200)

	cl, err := NewClient("server.com", "a_token")
	s.Assert().Nil(err)

	cl.Query("GET", "foo", nil)

	s.Assert().True(gock.IsDone())
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
