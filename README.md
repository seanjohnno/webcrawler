Currently working but still a work in progress*

## 1. Usage

### 1.1 Output to a directory

```
NewCrawlerBuilder("https://urltocrawl.com").
  BuildWithOutputDestination(/some/output/dir).
  Start()
```

### 1.2 Handle content yourself

```
NewCrawlerBuilder("https://urltocrawl.com").
  BuildWithOutputHandler(handler func(crawler Crawler, url string, content []byte) {
    // Handle
  }).
  Start()
```

### 1.3 With additional (optional) builder options

```
NewCrawlerBuilder("https://urltocrawl.com").
  WithMaxDepth(2).
  WithFilter(func(crawler Crawler, depth int, url string) bool {
    // return true to allow, false to ignore
  }).
  WithErrorHandler(func(crawler Crawler, err WebCrawlerError) {
	// Handle
  }).
  BuildWithOutputDestination(/some/output/dir).
  Start()
```

## *TODO
* Use a buffer for fileOutputHandler so its not loading all bytes into memory
* Adding additional links for ctrawler to fetch (use case = retrying errors)
* Link rewriting for resources not on same domain
* Add ability to use as command line utility as well as go lib
* Concurrency
