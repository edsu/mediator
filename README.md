# mediuminator

Look at Medium through the lens of Twitter.

mediuminator works by listening to a [filtered Twitter stream](https://dev.twitter.com/docs/api/1.1/post/statuses/filter) for `medium.com` urls. When tweets are found, the medium url is identified which is then used to extract relevant metadata from Medium pages, which is then saved off in the database.  Medium have done a nice job with their meta tags, so we could probably use an html5 parser to get at it.

## App

mediuminator can probably start out as a single page app that displays trending Medium stories. The stories will be listed by the number of times they've been mentioned on Twitter. Each story can use the author, collection, publication date and image url. The default view will be for trending stories in the last hour, but we should be able to add a control to switch to day or week.

It might also be fun to display what collections and users are trending, and how stories are being referenced from multiple collections.

## Data

Here's the data we should be able to store (object keys are in bold):

### Tweet

* **url**
* text
* story_url
* created
* user_url

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

* **url**
* name
* description
* image_url

### ShortUrl

* **short_url**
* long_url

While most object-store dbs will handle this sort of metadata (dynamodb,
mongo, redis, etc) a big part of this app is showing what's happening
over time, so do we need to think about time-series data differently?
