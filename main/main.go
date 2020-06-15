package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/kraczak/urlshort"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

var Global = map[string]string{}

func main() {
	mux := defaultMux()
	mapHandler := urlshort.MapHandler(Global, mux)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	return mux
}

func genURLHash(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)[:5]
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		fp := path.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, Global); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		formURL := r.FormValue("url")
		_, err := url.ParseRequestURI(formURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		urlHash := genURLHash(string(formURL))
		Global["/"+urlHash] = formURL
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
