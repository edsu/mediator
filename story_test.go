package mediator_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/edsu/mediator"
)

func TestStory(t *testing.T) {
	story := mediator.NewStory("https://medium.com/life-at-obvious/b7d6d31b4cd0")
	assert.Equal(t, story.Title, "A Medium-Sized Family — Life at Medium — Medium")
	assert.Equal(t, story.Description, "A team that works and plays well together.")
	assert.Equal(t, story.Url, "https://medium.com/jobs/b7d6d31b4cd0")
	assert.Equal(t, story.Author, "https://medium.com/@nmanekia5")
	assert.Equal(t, story.ImageUrl, "https://d262ilb51hltx0.cloudfront.net/max/800/1*pUkoGS6ur5wLBN06LkPjFw.jpeg")
}
