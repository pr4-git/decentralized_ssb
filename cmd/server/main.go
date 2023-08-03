package main

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"itsy/ssb/dba"
	securerpc "itsy/ssb/rpc"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

const appkey string = "ssbnet"

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

	//go RunServer(db, keypair)
	serverKey, err := base64.StdEncoding.DecodeString("1a99MIBEzGXc8gWbRIjRxgHRRdnf+KM7hinz4S0FIjE=")
	if err != nil {
		log.Fatalf("Server Key not available. %s", err)
	}
	go securerpc.SyncClient(db, keypair, ed25519.PublicKey(serverKey), appkey)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Blocking, press ctrl_c to exit...")
	<-done

}
