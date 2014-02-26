# mediator

[![Build Status](https://secure.travis-ci.org/edsu/mediator.png)](http://travis-ci.org/edsu/medinator)

mediator works by listening to a [filtered Twitter
stream](https://dev.twitter.com/docs/api/1.1/post/statuses/filter) for
`medium.com` urls. When tweets are found, the medium url is used to extract
relevant metadata from the Medium page, which is then saved off in the database.
Medium have done a nice job with their HTML metadata so this isn't as bad 
as it sounds.

## App

mediator is a single page app that displays trending Medium stories. The
stories are listed by the number of times they've been mentioned on Twitter. The
display of each story can use the author, collection, publication date and image
url. The default view will be for trending stories in the last hour, but a
control allows you to switch to daily and weekly views.

It might also be fun to display what collections and users are trending, and how 
stories are being referenced from multiple collections.

## Models

Here's the data we should be able to store (object keys are in bold):

### Tweet

* **url**
* text
* story_url
* created
* user_url (twitter.com url)

### Story

* **url**
* created
* title
* description
* image_url
* collection_url (optional)

### Collection

* **url**
* title

### User

* **url** (the medium.com url)
* name
* description
* image_url

### ShortUrl

* **short_url**
* long_url

## Time Series Data

TODO: need to come up with a scheme for modeling tweets and stories so that we
can report out top stories by the number of times they have been tweeted.


