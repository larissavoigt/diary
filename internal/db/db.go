package db

import (
	"database/sql"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

type User struct {
	ID    int64
	Token string
	Name  sql.NullString
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@/diary")
	if err != nil {
		panic(err)
	}
}

func CreateUser(token string) (string, error) {
	res, err := db.Exec("INSERT INTO users (token) VALUES(?)", token)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func FindUser(id string) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&u.ID, &u.Token, &u.Name)
	if err != nil {
		return nil, err
	} else {
		return u, nil
	}
}
