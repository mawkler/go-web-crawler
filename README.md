# Go Web Crawler

This is a toy web crawler built for the course [Learn Web Servers at Boot.dev](https://www.boot.dev/courses/build-web-crawler-golang).

## Usage

```sh
crawler <baseURL> <maxConcurrency> <maxPages>
```

Where `baseURL` is the URL to the website to crawl, `maxConcurrency` is the maximum number of go routines to spawn and `maxPages` is the maximum number web pages to crawl.
