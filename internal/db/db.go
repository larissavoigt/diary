package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

type fbres struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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
	api := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", token)
	r, err := http.Get(api)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	fb := &fbres{}
	err = json.Unmarshal(body, &fb)
	if err != nil {
		return "", err
	}

	res, err := db.Exec("INSERT INTO users (token, name) VALUES(?, ?)", token, fb.Name)
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
