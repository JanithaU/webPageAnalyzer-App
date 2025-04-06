package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

func isExternal(url, baseURL string) bool {
	// log.Println(url, baseURL)
	return !strings.HasPrefix(url, baseURL)
}

func analyzePage(url string) (*AnalysisResults, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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

	results := &AnalysisResults{
		HeadingCounts: make(map[string]int),
	}

	html5Regex := regexp.MustCompile(`(?i)<!DOCTYPE\s+html\s*>`)
	html4Regex := regexp.MustCompile(`(?i)<!DOCTYPE\s+HTML\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+([0-9]+\.[0-9]+)//EN"`)

	if html5Regex.MatchString(rawHTML) {
		results.HTMLVersion = "HTML5"
	} else if html4Regex.MatchString(rawHTML) {
		results.HTMLVersion = "HTML4"
	} else {
		results.HTMLVersion = "Unknown"
	}

	results.Url = url

	results.Title = doc.Find("title").Text()

	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		tagName := s.Get(0).Data
		if len(tagName) > 1 && tagName[0] == 'h' {
			level := tagName
			results.HeadingCounts[string(level)]++
		}
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

				if _, err := http.Get(href); err != nil {
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

	for results := range ch {
		log.Printf("URL: %v, External(url): %v, Inaccessible: %v  ", results.URL, results.isExternal, results.isInaccessible)

	}

	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, exists := s.Attr("action")
		if exists && (strings.Contains(action, "login") || strings.Contains(action, "signin")) {
			results.LoginFormPresent = true
		}
	})

	return results, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")

		log.Println("Fetching Started for Url:", url)
		results, err := analyzePage(url)
		if err != nil {
			// http.Error(w, "Error analyzing page: "+err.Error(), http.StatusInternalServerError)
			// return
			results = &AnalysisResults{
				ErrorMessage: "Error analyzing page: " + err.Error(),
			}
		}

		tmpl, err := template.ParseFiles("templates/result.html")
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			log.Println("Fetching Done for Url:", url)
		}

		err = tmpl.Execute(w, results)
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.ServeFile(w, r, "templates/index.html")
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
