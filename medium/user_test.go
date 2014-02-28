package medium_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/edsu/mediator/medium"
)

func TestUser(t *testing.T) {
	user := medium.GetUser("https://medium.com/@ChrisRosche/")
	assert.Equal(t, user.Url, "https://medium.com/@ChrisRosche")
	assert.Equal(t, user.Name, "Christopher Rosche")
	assert.Equal(t, user.Description, "\u00a0Recovering Renaissance man: writer, consultant, former congressional staffer, and   journalist trying to make sense of it all in my first novel.")
	assert.Equal(t, user.ImageUrl, "https://d262ilb51hltx0.cloudfront.net/max/800/0*DtYHWcsX17h2_4uq.png")
	assert.Equal(t, user.GoogleUrl, "https://plus.google.com/+ChristopherRosche")
	assert.Equal(t, user.TwitterUrl, "https://twitter.com/ChrisRosche")
}
