package db

import (
	"log"
	// "os"

	"github.com/jmoiron/sqlx"

	"github.com/elevin72/click-wars/internal/click"
	_ "github.com/mattn/go-sqlite3"
)

var onConnectSettings = [...]string{
	"pragma journal_mode = WAL;",
	"pragma synchronous = normal;",
	"pragma mmap_size = 30000000000;",
	"PRAGMA busy_timeout = 5000",
	"PRAGMA foreign_keys = true",
}

type Database struct {
	rw *sqlx.DB
	ro *sqlx.DB
}

func New() *Database {
	url := "click-wars.db?_journal=WAL&_timeout=5000&_fk=true&_txlock=immediate"

	rw, err := sqlx.Connect("sqlite3", url+"&mode=rw")
	if err != nil {
		log.Printf("Error creating rw connection: %v", err)
	}
	rw.SetMaxOpenConns(1)

	ro, err := sqlx.Connect("sqlite3", url+"&mode=ro")
	if err != nil {
		log.Printf("Error creating ro connection: %v", err)
	}
	// environment, exists := os.LookupEnv("CLICK_WARS_ENVIRONMENT") // TODO envs
	// if !exists {
	// 	environment = "dev"
	// }
	// if environment == "dev" {
	// 	rw.MustExec(`
	// 		DROP TABLE IF EXISTS clicks;
	// 		CREATE TABLE clicks (
	// 			id             INTEGER,
	// 			x              REAL,
	// 			y              REAL,
	// 			side           INTEGER ,
	// 			sender_uuid    TEXT,
	// 			line_position  INTEGER,
	// 			total_hits     INTEGER,
	// 			time           TEXT
	// 		) STRICT;
	// 	`)
	// }

	return &Database{
		rw: rw,
		ro: ro,
	}
}

func (d *Database) InsertClick(sc *click.ServerClick) {
	tx := d.rw.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO clicks (x, y, side, sender_uuid, line_position, total_hits, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		sc.X, sc.Y, sc.Side, sc.SenderUUID, sc.LinePosition, sc.TotalHits, sc.Time)
	if err != nil {
		log.Printf("error while inserting click: %v", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (d *Database) SelectClicks() click.ServerClick {
	c := click.ServerClick{}
	err := d.ro.Get(&c, `
		SELECT x, y, side, sender_uuid, line_position, total_hits, time
		FROM clicks
		WHERE time == (SELECT MAX(time) FROM clicks)
		LIMIT 1;`)
	if err != nil {
		log.Printf("error while selecting click: %v", err)
	}
	return c
}

func (d *Database) State() (linePosition int32, totalHits int32) {
	c := click.ServerClick{}
	err := d.ro.Get(&c, `
		SELECT x, y, side, sender_uuid, line_position, total_hits, time
		FROM clicks
		WHERE total_hits == (SELECT MAX(total_hits) FROM clicks);`)
	if err != nil {
		log.Printf("error while getting latest click: %v", err)
	}
	return c.LinePosition, c.TotalHits
}

func (d *Database) Close() {
	d.ro.Close()
	d.rw.Close()
}
