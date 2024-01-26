package main

import (
	"database/sql"
	"log"
	"m1thrandir225/your_time/util"

	_ "github.com/lib/pq"
)


func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	
}