package sqlite

import (
	"database/sql"
	 _ "modernc.org/sqlite"  // replace go-sqlite3
	"github.com/siddharth7actowiz/student_api/internal/config"
)

type SqliteStorage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err!=nil{
		return nil ,err
	}

	return &SqliteStorage{db: db}, nil

}



