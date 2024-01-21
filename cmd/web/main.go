package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/arian-nj/snippetbox/internals/models"
	_ "github.com/go-sql-driver/mysql"
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

type errorLevels struct {
	err  *log.Logger
	info *log.Logger
}

// dependency injection
type application struct {
	log      errorLevels
	snippets models.SnippetModel
}

func main() {
	app := application{}

	// Command line flags
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network Adrress")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySql data source name")
	flag.Parse()

	// setup logger
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.log.err = errorlog
	app.log.info = infolog

	// setup db
	db, err := OpenDB(*dsn)
	if err != nil {
		app.log.err.Println(err)
	} else {
		app.log.info.Println("Succesful connction to database")
	}

	defer db.Close()

	app.snippets = models.SnippetModel{DB: db} // look at me

	infolog.Printf("Starting server on http://127.0.0.1%s", cfg.addr)

	// set routers
	routes := app.routes()
	srv := http.Server{
		Addr:    cfg.addr,
		Handler: routes,
	}
	srv.ErrorLog = errorlog

	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
