package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ngnhub/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *mysql.SnippetModel
}

// Create info and error loggers
var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog *log.Logger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Defines argument variables
	addr := flag.String("port", ":4000", "Server address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Db url")
	// Parses arguments
	flag.Parse()

	//Open db connection
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v\n", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
