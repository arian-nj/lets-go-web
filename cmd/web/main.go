package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	// setup logger
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorlog: errorlog,
		infolog:  infolog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/cu", &home{})

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)

	infolog.Printf("Starting server on http://127.0.0.1%s", cfg.addr)

	srv := http.Server{Addr: cfg.addr, Handler: mux}
	srv.ErrorLog = errorlog
	err := srv.ListenAndServe()

	errorlog.Fatal(err)
}
