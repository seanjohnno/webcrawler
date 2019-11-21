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
* Error retry. Just adding links no good as you lose context (depth, parent)
	* Could also add a default retry handler (n times) 
* Use a buffer for fileOutputHandler so its not loading all bytes into memory
* Add ability to use as command line utility as well as go lib
	* Greedy grab
	* Pull/proxy content so you only pull down what you viewed for next time
	* Host content (localhost)
* Logging
* Concurrency
	* + addding a wait time / variable wait time. Makes sense to do this after
	adding concurreny/goroutines as would only have to re-do

 ## Edges
 * Fully qualified URL (with host) in link (still need to rewrite parent)

### Reminder
You're currently changing output handler to take http.Response so fileOutputHandler
is the thing which will do the url rewriting (normal handlers may not want urls rewritten
so it should happen in here)