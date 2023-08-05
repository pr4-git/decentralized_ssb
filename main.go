package main

import (
	"database/sql"
	"embed"
	"itsy/ssb/dba"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	logfile, err := os.OpenFile("ssb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening logfile. %v", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	keypair := GetIdentity()

	// init db
	sql.Register("sqlite3_with_crypto", &sqlite3.SQLiteDriver{
		Extensions: []string{
			"./sqlite_ext/crypto",
		},
	})
	db, err := sqlx.Open("sqlite3_with_crypto", "ssb.sqlite")
	if err != nil {
		log.Fatalf("Unable to open database. %s", err)
	}
	defer db.Close()
	dba.InitDB(db.DB)

	app := &Deps{
		db:      db,
		appkey:  "universal ssb network",
		keypair: *keypair,
	}

	go app.RunSyncServer()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "ssb",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
