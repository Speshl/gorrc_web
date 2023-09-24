package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Speshl/gorrc_web/internal/service/config"
	"github.com/Speshl/gorrc_web/internal/service/server"
	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

func main() {
	log.Println("starting server...")

	cfg := config.GetConfig()
	conn, err := getDBConn(cfg.DBCfg)
	if err != nil {
		panic(fmt.Errorf("error connecting to database: %w", err))
	}

	store := v1gorrc.NewStore(conn)

	server := server.NewServer(context.Background(), cfg.ServerCfg, store)

	server.StartServing()
}

func getDBConn(cfg config.DBCfg) (*dbr.Connection, error) {

	mysqlCfg := mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DBName:    cfg.DBName,
		ParseTime: true,
	}
	// Get a database handle.
	//return sql.Open("mysql", mysqlCfg.FormatDSN())
	return dbr.Open("mysql", mysqlCfg.FormatDSN(), nil)
}
