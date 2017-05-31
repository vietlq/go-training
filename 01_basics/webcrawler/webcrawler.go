package main

import (
    "fmt"
    "io"
    "os"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "sync"
    "golang.org/x/net/html"
    //"io/ioutil"
    "wcutil"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(curUrl string) (body string, urls []string, err error)
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

type ExtractedLinks struct {
    urls    []string
    imgs    []string
    scripts []string
    styles  []string
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

    body, urls, err := fetcher.Fetch(curUrl)
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

func ExtractAttr(z *html.Tokenizer, targetAttr string) string {
    key, val, moreAttr := z.TagAttr()
    attr := string(key)
    for len(attr) > 0 && attr != targetAttr && moreAttr {
        key, val, moreAttr = z.TagAttr()
        attr = string(key)
    }

    if attr == targetAttr {
        return strings.Trim(string(val), " \t\n")
    }

    return ""
}

func ExtractLinkCssHref(z *html.Tokenizer) string {
    attrs := make(map[string]string)

    key, val, moreAttr := z.TagAttr();
    for moreAttr {
        attrs[string(key)] = string(val)
        key, val, moreAttr = z.TagAttr();
    }
    attrs[string(key)] = string(val)

    if val, ok := attrs["rel"]; ok && val == "stylesheet" {
        if val, ok = attrs["href"]; ok {
            return strings.Trim(val, " \t\n")
        }
    }
    return ""
}

func ExtractLinks(refUrl string, r io.Reader) (resLinks ExtractedLinks) {
    z := html.NewTokenizer(r)
    parsedRefUrl, _ := url.Parse(refUrl)

    for {
        tt := z.Next()

        switch tt {
        case html.ErrorToken:
            return
        case html.StartTagToken:
            tn, ok := z.TagName()
            if !ok {
                continue
            }
            tagName := string(tn)
            switch tagName {
            case "a":
                // Extract HREF from A
                curUrl := ExtractAttr(z, "href")
                curUrl = wcutil.NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.urls = append(resLinks.urls, curUrl)
                }
            case "img":
                // Extract SRC from IMG
                curUrl := ExtractAttr(z, "src")
                curUrl = wcutil.NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.imgs = append(resLinks.imgs, curUrl)
                }
            case "script":
                // Extract SRC from SCRIPT
                curUrl := ExtractAttr(z, "src")
                curUrl = wcutil.NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.scripts = append(resLinks.scripts, curUrl)
                }
            case "link":
                // Extract CSS HREF from LINK with REL="stylesheet"
                curUrl := ExtractLinkCssHref(z)
                curUrl = wcutil.NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.styles = append(resLinks.styles, curUrl)
                }
            }
        }
    }

    return
}

func (f *PageFetcher) Fetch(curUrl string) (string, []string, error) {
    // Use the cache
    if res, ok := f.visited[curUrl]; ok {
        return res.body, res.urls, nil
    }
    // Fetch if not in cache
    resp, err := http.Get(curUrl)
    // Report error
    if err != nil {
        return "", nil, fmt.Errorf("not found: %s", curUrl)
    }
    // Read the response body
    defer resp.Body.Close()
    // Extract URLs
    resLinks := ExtractLinks(curUrl, resp.Body)
    fmt.Println(resLinks.imgs, resLinks.scripts)

    return "string(body)", resLinks.urls, nil
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
