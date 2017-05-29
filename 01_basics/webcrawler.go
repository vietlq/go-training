package main

import (
    "fmt"
    "sync"
    "os"
    "net/http"
    "io"
    "strconv"
    //"io/ioutil"
    "golang.org/x/net/html"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}
// PageFetcher is Fetcher that returns canned results.
type PageFetcher struct {
    visited map[string]*FetchResult
}

type FetchResult struct {
    body string
    urls []string
}

func ExtractHref(z *html.Tokenizer) string {
    key, val, moreAttr := z.TagAttr()
    attr := string(key)
    for len(attr) > 0 && attr != "href" && moreAttr {
        key, val, moreAttr = z.TagAttr()
        attr = string(key)
    }

    if attr == "href" {
        return string(val)
    }

    return ""
}

func ExtractLinks(r io.Reader) []string {
    z := html.NewTokenizer(r)
    urls := []string{}

    for {
        tt := z.Next()

        switch tt {
        case html.ErrorToken:
            return urls
        case html.StartTagToken:
            tn, _ := z.TagName()
            // Extract HREF from A
            if len(tn) == 1 && tn[0] == 'a' {
                url := ExtractHref(z)
                if (len(url) > 0) {
                    urls = append(urls, url)
                }
            }
        }
    }

    return urls
}

func (f *PageFetcher) Fetch(url string) (string, []string, error) {
    // Use the cache
    if res, ok := f.visited[url]; ok {
        return res.body, res.urls, nil
    }
    // Fetch if not in cache
    resp, err := http.Get(url)
    // Report error
    if err != nil {
        return "", nil, fmt.Errorf("not found: %s", url)
    }
    // Read the response body
    defer resp.Body.Close()
    // Extract URLs
    urls := ExtractLinks(resp.Body)

    return "string(body)", urls, nil
}

type VisitDict struct {
    visits map[string]bool
    mux    sync.RWMutex
}

func (vd *VisitDict) Visited(url string) bool {
    vd.mux.RLock()
    defer vd.mux.RUnlock()
    // https://blog.golang.org/go-maps-in-action
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
    ch chan FetchResult, wg *sync.WaitGroup) {
    // https://stackoverflow.com/questions/19892732/all-goroutines-are-asleep-deadlock
    defer wg.Done()

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

    fmt.Printf("found: %s\n", url)
    ch <- FetchResult{body, urls}

    for _, u := range urls {
        wg.Add(1)
        go Crawl(u, depth-1, fetcher, visitDict, ch, wg)
    }
    return
}

func monitorWorker(wg *sync.WaitGroup, ch chan FetchResult) {
    wg.Wait()
    close(ch)
}

func implementation(depth int, seeds []string) {
    fetcher := PageFetcher{visited: make(map[string]*FetchResult)}
    visitDict := VisitDict{visits: make(map[string]bool)}
    ch := make(chan FetchResult)
    wg := &sync.WaitGroup{}

    // Launch the workers based on seed URLs
    for _, url := range seeds {
        wg.Add(1)
        go Crawl(url, depth, &fetcher, &visitDict, ch, wg)
    }

    // Monitor the workers
    go monitorWorker(wg, ch)

    // Wait and reap results
    for {
        select {
        case res, ok := <-ch:
            if !ok {
                return
            }
            fmt.Println("Got result with URLs:", len(res.urls))
        }
    }
}

func UsageExit() {
    fmt.Println("Usage: Program Depth <'URLs'> <'separated'> <'by'> <'space'>")
    os.Exit(-1)
}

func main() {
    if (len(os.Args) < 3) {
        UsageExit()
    }
    depth, err := strconv.Atoi(os.Args[1])
    if err != nil {
        UsageExit()
    }
    implementation(depth, os.Args[2:])
}
