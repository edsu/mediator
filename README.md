# medium-inator

medium-inator is an app that looks at Medium through the lens of Twitter.

## Entities

medium-inator works by listening to a [filtered Twitter
stream](https://dev.twitter.com/docs/api/1.1/post/statuses/filter) for
`medium.com`. When tweets are found, the medium.url is identified (possibly via
link-shortening) and is used to extract relevant metadata from Medium pages. 
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
