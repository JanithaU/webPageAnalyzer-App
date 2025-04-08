package analysis

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type AnalysisResults struct {
	Url               string
	HTMLVersion       string
	Title             string
	HeadingCounts     map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	LoginFormPresent  bool
	ErrorMessage      string
}

func isExternal(link, base string) bool {
	baseURL, err := url.Parse(base)
	if err != nil {
		return false
	}
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	resolvedURL := baseURL.ResolveReference(linkURL)
	return resolvedURL.Host != baseURL.Host
}

func AnalyzePage(url string) (*AnalysisResults, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d for URL %s", res.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	rawHTML, err := doc.Html()
	if err != nil {
		return nil, err
	}

	results := &AnalysisResults{HeadingCounts: make(map[string]int)}
	results.Url = url

	// Determine HTML version
	htmlVersionPatterns := map[string]*regexp.Regexp{
		"HTML5":                  regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s*>`),
		"HTML 4.01 Strict":       regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+4\.01//EN"\s+"http://www\.w3\.org/TR/html4/strict\.dtd">`),
		"HTML 4.01 Transitional": regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+4\.01\s+Transitional//EN"\s+"http://www\.w3\.org/TR/html4/loose\.dtd">`),
		"HTML 4.01 Frameset":     regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+4\.01\s+Frameset//EN"\s+"http://www\.w3\.org/TR/html4/frameset\.dtd">`),
		"HTML 3.2":               regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+3\.2\s+Final//EN">`),
		"XHTML 1.0 Strict":       regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD\s+XHTML\s+1\.0\s+Strict//EN"\s+"http://www\.w3\.org/TR/xhtml1/DTD/xhtml1-strict\.dtd">`),
		"XHTML 1.0 Transitional": regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD\s+XHTML\s+1\.0\s+Transitional//EN"\s+"http://www\.w3\.org/TR/xhtml1/DTD/xhtml1-transitional\.dtd">`),
		"XHTML 1.0 Frameset":     regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD\s+XHTML\s+1\.0\s+Frameset//EN"\s+"http://www\.w3\.org/TR/xhtml1/DTD/xhtml1-frameset\.dtd">`),
		"XHTML 1.1":              regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s+PUBLIC\s+"-//W3C//DTD\s+XHTML\s+1\.1//EN"\s+"http://www\.w3\.org/TR/xhtml11/DTD/xhtml11\.dtd">`),
	}

	detectedVersion := "Unknown"
	for version, pattern := range htmlVersionPatterns {
		if pattern.MatchString(rawHTML) {
			detectedVersion = version
			break
		}
	}
	results.HTMLVersion = detectedVersion

	results.Title = doc.Find("title").Text()

	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		tagName := s.Get(0).Data
		results.HeadingCounts[tagName]++
	})

	baseURL := res.Request.URL.String()
	ch := make(chan struct {
		URL            string
		isExternal     bool
		isInaccessible bool
	}, 10)
	var wg sync.WaitGroup

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			wg.Add(1)
			go func(href string) {
				defer wg.Done()
				isExt := isExternal(href, baseURL)
				if isExt {
					results.ExternalLinks++
				} else {
					results.InternalLinks++
				}
				_, err := http.Get(href)
				if err != nil {
					results.InaccessibleLinks++
				}
				ch <- struct {
					URL            string
					isExternal     bool
					isInaccessible bool
				}{href, isExt, err != nil}
			}(href)
		}
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	for r := range ch {
		log.Infof("URL: %v, External: %v, Inaccessible: %v", r.URL, r.isExternal, r.isInaccessible)
	}

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		for _, attr := range []string{"action", "href", "id", "class"} {
			if value, exists := s.Attr(attr); exists {
				if strings.Contains(value, "login") || strings.Contains(value, "signin") {
					results.LoginFormPresent = true
					return
				}
			}
		}
	})

	return results, nil
}
