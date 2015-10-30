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

type fbUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID         int64
	FacebookID string
	Token      string
	Name       sql.NullString
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

func FindUser(id string) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(
		&u.ID, &u.FacebookID, &u.Token, &u.Name)
	if err != nil {
		return nil, err
	} else {
		return u, nil
	}
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
