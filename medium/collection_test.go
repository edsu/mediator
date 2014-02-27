package medium_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/edsu/mediator/medium"
)

func TestCollection(t *testing.T) {
	coll := medium.GetCollection("https://medium.com/life-at-obvious/")
	assert.Equal(t, coll.Title, "Life at Medium — Medium")
	assert.Equal(t, coll.Url, "https://medium.com/life-at-obvious")
	assert.Equal(t, coll.Description, "A unique experience")
	assert.Equal(t, coll.ImageUrl, "https://d262ilb51hltx0.cloudfront.net/max/800/0*uPsGAjmMo7Yr49AF.jpeg")
}
