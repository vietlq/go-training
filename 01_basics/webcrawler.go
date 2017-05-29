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

func ExtractAttr(z *html.Tokenizer, targetAttr string) string {
    key, val, moreAttr := z.TagAttr()
    attr := string(key)
    for len(attr) > 0 && attr != targetAttr && moreAttr {
        key, val, moreAttr = z.TagAttr()
        attr = string(key)
    }

    if attr == targetAttr {
        return string(val)
    }

    return ""
}

func ExtractLinkCssHref(z *html.Tokenizer) string {
    attrs = make(map[string]string)

    key, val, moreAttr := z.TagAttr();
    for moreAttr {
        attrs[string(key)] = string(val)
        key, val, moreAttr = z.TagAttr();
    }
    attrs[string(key)] = string(val)

    if val, ok := attrs["rel"]; ok && val == "stylesheet" {
        if val, ok = attrs["href"]; ok {
            return val
        }
        return ""
    }
    return ""
}

func ExtractLinks(r io.Reader) (urls []string, imgs []string, scripts []string) {
    z := html.NewTokenizer(r)

    for {
        tt := z.Next()

        switch tt {
        case html.ErrorToken:
            return urls, imgs, scripts
        case html.StartTagToken:
            tn, ok := z.TagName()
            if !ok {
                continue
            }
            tagName := string(tn)
            switch tagName {
            case "a":
                // Extract HREF from A
                url := ExtractAttr(z, "href")
                if (len(url) > 0) {
                    urls = append(urls, url)
                }
            case "img":
                // Extract SRC from IMG
                url := ExtractAttr(z, "src")
                if (len(url) > 0) {
                    imgs = append(imgs, url)
                }
            case "script":
                // Extract SRC from SCRIPT
                url := ExtractAttr(z, "src")
                if (len(url) > 0) {
                    scripts = append(scripts, url)
                }
            case "link":
                // Extract CSS HREF from LINK with REL="stylesheet"
            }
        }
    }

    return urls, imgs, scripts
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
    urls, imgs, _ := ExtractLinks(resp.Body)
    fmt.Println(imgs)

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
