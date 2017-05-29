package normalise_url

import (
    "testing"
    "net/url"
)

type TestCase struct {
    refUrl *url.URL
    curUrl string
    want string
}

func CaseNoSchemeNoHost(tests []TestCase) {
    refUrlStr1 := "http://code2.pro/cool/course"
    refUrlStr2 := "http://code2.pro/cool/course/"
    curUrl := "iscoming/soon"
    want1 := "http://code2.pro/cool/iscoming/soon"
    want2 := "http://code2.pro/cool/course/iscoming/soon"
    refUrl1, _ := url.Parse(refUrlStr1)
    refUrl2, _ := url.Parse(refUrlStr2)
    tests = append(tests, TestCase{refUrl1, curUrl, want1})
    tests = append(tests, TestCase{refUrl2, curUrl, want2})
}

func Test(t *testing.T) {
    tests := make([]TestCase, 10)
    for _, c := range tests {
        got := NormaliseUrl(c.refUrl, c.curUrl)
        if got != c.want {
            t.Errorf("NormaliseUrl(%q, %q) = %q, expected %q", c.refUrl, c.curUrl, got, c.want)
        }
    }
}
