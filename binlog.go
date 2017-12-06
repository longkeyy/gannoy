package gannoy

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type BinLog struct {
	Path string
	db   *sql.DB
	stms *sql.Stmt
}

func NewBinLog(path string) BinLog {
	return BinLog{Path: path}
}

func (b *BinLog) Open() error {
	db, err := sql.Open("sqlite3", b.Path)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`create table if not exists "features" ("key" integer primary key, "action" integer, "features" blob)`,
	)
	if err != nil {
		return err
	}
	stms, err := db.Prepare(
		`insert or replace into features values (?, ?, ?)`,
	)
	if err != nil {
		return err
	}
	b.db = db
	b.stms = stms
	return nil
}

func (b BinLog) Close() error {
	err := b.stms.Close()
	if err != nil {
		return err
	}
	return b.db.Close()
}

func (b BinLog) Add(key, action int, features []byte) error {
	_, err := b.stms.Exec(key, action, features)
	return err
}
