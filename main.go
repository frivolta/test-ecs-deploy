package main

import (
	"birdie/api"
	db "birdie/db/sqlc"
	"birdie/util"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".", "./serviceAccountKey.json")
	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	fmt.Println(gin.EnvGinMode)
	fmt.Println("DB Address", config.DBSource)
	if err != nil {
		log.Fatal("Cannot load configuration", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	err = server.Start(config.ServerAddress)

}
