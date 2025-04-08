package handler

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/JanithaU/webPageAnalyzer-App/internal/analysis"

	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Url analyzer handler
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the absolute path to your template file

	if r.Method == http.MethodPost {
		r.ParseForm()
		url := r.FormValue("url")

		log.Info("Fetching Started for Url:", url)
		results, err := analysis.AnalyzePage(url)
		if err != nil {
			log.Error("Error analyzing page: " + err.Error())
			results = &analysis.AnalysisResults{
				ErrorMessage: "Error analyzing page: " + err.Error(),
			}
		}
		templatePath := filepath.Join(cwd, "web", "templates", "result.html")
		// tmpl, err := template.ParseFiles("../../web/templates/result.html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Error("Error parsing template: ", err)
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}

		log.Info("Fetching Done for Url:", url)

		if err := tmpl.Execute(w, results); err != nil {
			log.Error("Error rendering template: ", err)
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	} else {
		// log.Info("sdf")
		templatePath := filepath.Join(cwd, "web", "templates", "index.html")
		// http.ServeFile(w, r, "../../web/templates/index.html")
		http.ServeFile(w, r, templatePath)

	}
}

// 404 handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNotFound)
	templatePath := filepath.Join(cwd, "web", "templates", "404.html")
	// tmpl, err := template.ParseFiles("../../web/templates/404.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Error("Error rendering template(404): ", err)
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, nil)
}
