package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/JanithaU/webPageAnalyzer-App/internal/handler"
)

// RegisterRoutes sets up the application's route handlers.
func RegisterRoutes() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	statPathJoin := filepath.Join(cwd, "web", "static")
	// staticPath := filepath.Join("..", "..", "web", "static")
	staticPath := filepath.Join(statPathJoin)
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/404", handler.NotFoundHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handler.NotFoundHandler(w, r)
			return
		}
		handler.AnalyzeHandler(w, r)
	})
}
