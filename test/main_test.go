package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var (
	//db *sql.DB
	//err error
)

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", "root:Hannamysql.1518@tcp(127.0.0.1:49547)/otabe")
	if err != nil {
		log.Fatalf("Error validating sql.Open arguments %v", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error verifying with db.Ping %v", err)
		panic(err)
	}
	os.Exit(m.Run())


}

