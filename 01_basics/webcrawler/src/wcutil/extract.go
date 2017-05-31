package wcutil

import (
    "io"
    "net/url"
    "strings"
    "golang.org/x/net/html"
)

type ExtractedLinks struct {
    // Properties and methods starting with an UPPER case are public/exported
    // Properties and methods starting with a lower case are private/unexported
    Urls    []string
    Imgs    []string
    Scripts []string
    Styles  []string
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
                curUrl = NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.Urls = append(resLinks.Urls, curUrl)
                }
            case "img":
                // Extract SRC from IMG
                curUrl := ExtractAttr(z, "src")
                curUrl = NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.Imgs = append(resLinks.Imgs, curUrl)
                }
            case "script":
                // Extract SRC from SCRIPT
                curUrl := ExtractAttr(z, "src")
                curUrl = NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.Scripts = append(resLinks.Scripts, curUrl)
                }
            case "link":
                // Extract CSS HREF from LINK with REL="stylesheet"
                curUrl := ExtractLinkCssHref(z)
                curUrl = NormaliseUrl(parsedRefUrl, curUrl)
                if (len(curUrl) > 0) {
                    resLinks.Styles = append(resLinks.Styles, curUrl)
                }
            }
        }
    }

    return
}
