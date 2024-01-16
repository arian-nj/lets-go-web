package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type home struct{}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is my byte"))
}

type Config struct {
	addr      string
	staticDir string
}

var cfg Config

func main() {
	// Get cli glags
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network Adrress")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	fmt.Println("Starting")

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/cu", &home{})

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/snippet/view", SnippetView)
	mux.HandleFunc("/snippet/create", SnippetCreate)

	log.Print("Starting server on http://127.0.0.1" + cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
}
