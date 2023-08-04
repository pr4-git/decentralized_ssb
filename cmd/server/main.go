package main

import (
	"database/sql"
	"flag"
	"itsy/ssb/dba"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

const appkey string = "ssbnet"

var httpPort *string

func init() {
	httpPort = flag.String("port", ":8080", "-port=:8080")
	flag.Parse()
}

func main() {
	// Create or get your identity
	keypair := GetIdentity()

	sql.Register("sqlite3_with_crypto", &sqlite3.SQLiteDriver{
		Extensions: []string{
			"./sqlite_ext/crypto",
		},
	})

	db, err := sqlx.Open("sqlite3_with_crypto", "./foo.db")
	if err != nil {
		log.Fatalln("Unable to open database. Cannot recover from this error.")
	}
	defer db.Close()
	dba.InitDB(db.DB)

	app := &Deps{
		db:      db,
		appkey:  appkey,
		keypair: *keypair,
	}

	go app.RunSyncServer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/api", func() http.Handler {
		r := chi.NewRouter()
		setupRoutes(r, app)
		return r
	}())

	http.ListenAndServe(*httpPort, r)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Blocking, press ctrl_c to exit...")
	<-done

}
