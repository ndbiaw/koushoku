package database

import (
	_ "embed"
	"sync"

	"database/sql"
	"fmt"
	"log"

	. "koushoku/config"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var Conn *sql.DB

//go:embed schema.sql
var schema []byte
var once sync.Once

func Init() {
	once.Do(func() {
		cfg := Config.Database
		dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Passwd, cfg.SSLMode)

		conn, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Fatalln(err)
		}

		if err := conn.Ping(); err != nil {
			log.Fatalln(err)
		}

		if _, err = conn.Exec(string(schema)); err != nil && err != sql.ErrNoRows {
			log.Fatalln(err)
		}

		Conn = conn
		boil.SetDB(conn)
	})
}
