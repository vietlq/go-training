package main

import (
    "fmt"
    "os"
    "net/http"
    "strconv"
    "sync"
    //"io/ioutil"
    "wcutil"
    "strings"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(curUrl string, allowedStatus map[int]struct{}) (body string, urls []string, err error)
}
// PageFetcher is Fetcher that returns canned results.
type PageFetcher struct {
    visited map[string]*FetchResult
}

type FetchResult struct {
    body string
    urls []string
}

type VisitDict struct {
    visits map[string]bool
    mux    sync.RWMutex
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

func UsageExit() {
    fmt.Println("Usage: Program Depth <'URLs'> <'separated'> <'by'> <'space'>")
    os.Exit(-1)
}

func implementation(depth int, seeds []string) {
    fetcher := PageFetcher{visited: make(map[string]*FetchResult)}
    visitDict := VisitDict{visits: make(map[string]bool)}
    ch := make(chan FetchResult)
    wg := &sync.WaitGroup{}

    // Launch the workers based on seed URLs
    for _, curUrl := range seeds {
        wg.Add(1)
        go Crawl(curUrl, depth, &fetcher, &visitDict, ch, wg)
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

    fmt.Println("Total links visited:", len(visitDict.visits))
}

// Crawl uses fetcher to recursively crawl
// pages starting with curUrl, to a maximum of depth.
func Crawl(curUrl string, depth int,
    fetcher Fetcher, visitDict *VisitDict,
    ch chan FetchResult, wg *sync.WaitGroup) {
    // https://stackoverflow.com/questions/19892732/all-goroutines-are-asleep-deadlock
    defer wg.Done()

    if depth <= 0 {
        return
    }

    // Don't visit if already done so
    if visitDict.Visited(curUrl) {
        return
    }
    visitDict.Visit(curUrl)

    allowedStatus := make(map[int]struct{})
    body, urls, err := fetcher.Fetch(curUrl, allowedStatus)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("found: %s\n", curUrl)
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

func (f *PageFetcher) Fetch(curUrl string, allowedStatus map[int]struct{}) (string, []string, error) {
    // Use the cache
    if res, ok := f.visited[curUrl]; ok {
        return res.body, res.urls, nil
    }

    // Fetch if not in cache
    resp, err := http.Get(curUrl)

    // Report error
    if err != nil {
        return "", nil, fmt.Errorf("URL not found: %s", curUrl)
    }

    // If no specific HTTP StatusCode, accept all
    if len(allowedStatus) > 0 {
        // Check if the return code is allowed
        if _, ok := allowedStatus[resp.StatusCode]; !ok {
            return "", nil, fmt.Errorf("Bad HTTP Status: %d returned by URL %s", resp.StatusCode, curUrl)
        }
    }
    fmt.Printf("StatusCode: %d, URL: %s\n", resp.StatusCode, curUrl)
    fmt.Printf("Proto: %s, URL: %s\n", resp.Proto, curUrl)

    header := resp.Header
    if err = CheckFetchErrors(curUrl, header); err != nil {
        return "", nil, err
    }

    // Read the response body
    defer resp.Body.Close()
    // Extract URLs
    resLinks := wcutil.ExtractLinks(curUrl, resp.Body)
    fmt.Println(resLinks.Imgs, resLinks.Scripts)

    return "string(body)", resLinks.Urls, nil
}

func CheckFetchErrors(curUrl string, header http.Header) error {
    // Check if it's really HTML before trying to extract anything
    val, ok := header["Content-Type"]
    if !ok || len(val) != 1 {
        return fmt.Errorf("Bad Content-Type, expected HTML: %q returned by URL %s", val, curUrl)
    }
    contentType := strings.ToLower(val[0])
    cntTypeParts := []string{}
    for _, v := range strings.Split(contentType, ";") {
        cntTypeParts = append(cntTypeParts, strings.TrimSpace(v))
    }
    if cntTypeParts[0] != "text/html" {
        return fmt.Errorf("Bad Content-Type, expected HTML: %q returned by URL %s", cntTypeParts[0], curUrl)
    }
    fmt.Printf("Content-Type: <%s>, URL: %s\n", contentType, curUrl)

    // Check the Content-Length
    val, ok = header["Content-Length"]
    const MAX_CONTENT_LEN = 1024*1024
    if ok {
        contentLen, convErr := strconv.Atoi(val[0])
        if convErr == nil && contentLen > MAX_CONTENT_LEN {
            return fmt.Errorf("Bad Content-Length, must not exceed %v (bytes), actual %v (bytes) returned by URL %s", MAX_CONTENT_LEN, val, curUrl)
        }
    }

    return nil
}

func (vd *VisitDict) Visited(curUrl string) bool {
    vd.mux.RLock()
    defer vd.mux.RUnlock()
    // https://blog.golang.org/go-maps-in-action
    _, ok := vd.visits[curUrl]
    return ok
}

func (vd *VisitDict) Visit(curUrl string) {
    vd.mux.Lock()
    vd.visits[curUrl] = true
    vd.mux.Unlock()
}
