# mediator

[![Build Status](https://secure.travis-ci.org/edsu/mediator.png)](http://travis-ci.org/edsu/mediator)

mediator works by listening to a [filtered Twitter
stream](https://dev.twitter.com/docs/api/1.1/post/statuses/filter) for
`medium.com` urls. When tweets are found, the medium url is used to extract
story, author and collection metadata from the Medium page, which is then 
saved off in the database. Medium have done a nice job with their HTML 
metadata so this isn't as bad as it sounds.

## App

mediator is a single page app that displays trending Medium stories. The
stories are listed by the number of times they've been mentioned on Twitter. The
display of each story can use the author, collection, publication date and image
url. The default view will be for trending stories in the last hour, but a
control allows you to switch to daily and weekly views.

It might also be fun to display what collections and users are trending, and 
how stories are being referenced from multiple collections. And, of course, 
it would be nice to display the updates in realtime. :cake:

## Models

Here's the data we should be able to store (object keys are in bold):

### Tweet

* **url**
* text
* created
* story (Story)
* author (TwitterUser)

### TwitterUser

* **url**
* name
* avatar_url

### Story

* **url**
* created
* title
* description
* image_url
* author (MediumUser)
* collection (Collection)

### MediumUser

* **url**
* name
* description
* image_url

### Collection

* **url**
* title
* description

## Time Series Data

TODO: need to come up with a scheme for modeling tweets and stories so that we
can efficiently report out top stories by the number of times they have been tweeted.
