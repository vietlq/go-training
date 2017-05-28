package main

import (
    "fmt"
    "sync"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

type VisitDict struct {
    visits map[string]bool
    mux    sync.Mutex
}

func (vd *VisitDict) Visited(url string) bool {
    vd.mux.Lock()
    defer vd.mux.Unlock()
    _, ok := vd.visits[url]
    return ok
}

func (vd *VisitDict) Visit(url string) {
    vd.mux.Lock()
    vd.visits[url] = true
    vd.mux.Unlock()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int,
    fetcher Fetcher, visitDict *VisitDict,
    ch chan fakeResult, wg *sync.WaitGroup) {
    // https://stackoverflow.com/questions/19892732/all-goroutines-are-asleep-deadlock
    defer wg.Done()

    // TODO: Fetch URLs in parallel.
    // TODO: Don't fetch the same URL twice.
    // This implementation doesn't do either:
    if depth <= 0 {
        return
    }

    // Don't visit if already done so
    if visitDict.Visited(url) {
        return
    }
    visitDict.Visit(url)

    body, urls, err := fetcher.Fetch(url)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("found: %s %q\n", url, body)
    ch <- fakeResult{body, urls}

    for _, u := range urls {
        wg.Add(1)
        go Crawl(u, depth-1, fetcher, visitDict, ch, wg)
    }
    return
}

func monitorWorker(wg *sync.WaitGroup, ch chan fakeResult) {
    wg.Wait()
    close(ch)
}

func implementation() {
    visitDict := VisitDict{visits: make(map[string]bool)}
    ch := make(chan fakeResult)
    wg := &sync.WaitGroup{}

    wg.Add(1)
    go Crawl("http://golang.org/", 4, fetcher, &visitDict, ch, wg)
    go monitorWorker(wg, ch)

    for {
        select {
        case res, ok := <-ch:
            if !ok {
                return
            }
            fmt.Println("res =", res)
        }
    }
}

func main() {
    implementation()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := f[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}
