package database

import "github.com/ostafen/clover"

type Databaser interface {
	RetrieveAndDelete(topic string) (Backups, error)
	Store(data Backup) error
}

type DatabaseDriver string

const (
	MYSQL  DatabaseDriver = "mysql"
	CLOVER DatabaseDriver = "clover"
)

type Database struct {
	DatabaseDriver DatabaseDriver
}

func (d *Database) NewCloverDatabase(db *clover.DB) Databaser {
	return NewCloverDatabase(db)
}
