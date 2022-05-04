package main

//DONT FORGET _ "github.com/lib/pq" ELSE FAILURE
import (
	"database/sql"
	"log"

	"github.com/kingsleyocran/simple_bank_bankend/api"
	db "github.com/kingsleyocran/simple_bank_bankend/db/sqlc"
	"github.com/kingsleyocran/simple_bank_bankend/util"
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

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
