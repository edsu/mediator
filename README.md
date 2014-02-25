# medium-inator

Look at Medium through the lens of Twitter.

## Entities

medium-inator works by listening to a [filtered Twitter
stream](https://dev.twitter.com/docs/api/1.1/post/statuses/filter) for
`medium.com`. When tweets are found, the medium url is identified (possibly via
link-shortening) and is used to extract relevant metadata from Medium pages.
Medium have done a nice job with their meta tags, so this shouldn't be too
tough.

Here's the data we can store.

### Tweet

* **tweet_id**
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
* collection_url

### Collection

* **url**
* title

### User

* name
* description
* image_url

## ShortUrl

* **short_url**
* long_url

## Questions

While most object-store dbs will handle this sort of metadata (dynamodb,
mongo, redis, etc) a big part of this app is showing what's happening
over time, so do we need to think about time-series data differently?
