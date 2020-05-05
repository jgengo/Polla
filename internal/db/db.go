package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Init() error {
	database, _ = sql.Open("sqlite3", "./database.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS poll (id INTEGER PRIMARY KEY, message TEXT, ts TEXT, channel TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS answer (id INTEGER PRIMARY KEY, poll_id INTEGER, message TEXT)")
	statement.Exec()
	return nil
}

func GenerateText(TS string) {
	var message string
	err := database.QueryRow("select message from poll where ts = ?", TS).Scan(&message)
	if err != nil {
		fmt.Printf("failed to get poll: %v\n", err)
	}
	fmt.Printf("message: %s", message)
}

func AddPoll(message, channel string) (int64, error) {
	statement, _ := database.Prepare("INSERT INTO poll (message, channel) VALUES (?, ?)")
	res, _ := statement.Exec(message, channel)
	lastID, _ := res.LastInsertId()
	fmt.Printf("=== lastID: %d", lastID)
	return lastID, nil
}

func UpdatePollTS(id int64, ts string) {
	statement, _ := database.Prepare("update poll set ts=? where id=?")
	statement.Exec(ts, id)
}
