package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/larissavoigt/diary/internal/entry"
	"github.com/larissavoigt/diary/internal/user"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

type fbUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@/diary")
	if err != nil {
		panic(err)
	}
}

func CreateUser(token string) (string, error) {
	user, err := getFBInfo(token)
	if err != nil {
		return "", err
	}
	res, err := db.Exec(`INSERT INTO users (facebook_id, token, name)
	VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE
	token=VALUES(token), name=VALUES(name)`, user.ID, token, user.Name)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func FindUser(id string) (*user.User, error) {
	u := &user.User{}
	err := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(
		&u.ID, &u.FacebookID, &u.Token, &u.Name)
	if err != nil {
		return nil, err
	} else {
		return u, nil
	}
}

func CreateEntry(id int64, rate, desc string) (string, error) {
	n, err := strconv.Atoi(rate)
	if err != nil {
		return "", err
	}
	res, err := db.Exec(`INSERT INTO entries (user_id, rate, description)
	VALUES(?, ?, ?)`, id, n, desc)
	if err != nil {
		return "", err
	}
	e, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(e, 10), nil
}

func FindUserEntries(id int64) ([]entry.Entry, error) {
	var entries []entry.Entry
	rows, err := db.Query("select id, rate, description from entries where user_id = ? ORDER BY id DESC LIMIT 10", id)
	if err != nil {
		return entries, err
	}
	defer rows.Close()
	for rows.Next() {
		e := entry.Entry{}
		err := rows.Scan(&e.ID, &e.Rate, &e.Description)
		if err != nil {
			return entries, err
		}
		entries = append(entries, e)
	}
	err = rows.Err()
	return entries, err
}

func getFBInfo(token string) (*fbUser, error) {
	api := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", token)
	r, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	u := &fbUser{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
