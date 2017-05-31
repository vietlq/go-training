package wcutil

import (
    "fmt"
    "net/url"
    "strings"
)

func NormaliseUrl(refUrl *url.URL, curUrl string) string {
    if curUrl == "" {
        return ""
    }

    // Infer HTTP or HTTPS
    if len(curUrl) >= 2 && curUrl[:2] == "//" {
        return fmt.Sprintf("%s:%s", refUrl.Scheme, curUrl)
    }

    // Use the front part of refUrl
    if curUrl[0] == '/' {
        return fmt.Sprintf("%s://%s%s", refUrl.Scheme, refUrl.Host, curUrl)
    }

    // Make sure we have valid URL
    u, err := url.Parse(curUrl)
    if err != nil {
        return ""
    }

    // Handle relative URLs
    if u.Scheme == "" {
        tpath := refUrl.Path
        lidx := strings.LastIndex(tpath, "/")
        return fmt.Sprintf("%s://%s%s/%s", refUrl.Scheme, refUrl.Host, tpath[:lidx], u.Path)
    }

    return curUrl
}
