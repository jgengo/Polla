package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

// Init opens and  creates the database if needed
func Init() error {
	database, _ = sql.Open("sqlite3", "./database.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS poll (id INTEGER PRIMARY KEY, message TEXT, ts TEXT, channel TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS answer (id INTEGER PRIMARY KEY, poll_id INTEGER, message TEXT, author TEXT)")
	statement.Exec()
	return nil
}

// GenerateText will generate the text for the posted message
func GenerateText(pollID int64) string {
	var message string
	err := database.QueryRow("select message from poll where id = ?", pollID).Scan(&message)
	if err != nil {
		fmt.Printf("failed to get poll: %v\n", err)
	}

	message += "\n*Responses:*\n"

	var tmp string
	rows, _ := database.Query("select message from answer where poll_id = ?", pollID)
	for rows.Next() {
		rows.Scan(&tmp)
		message += fmt.Sprintf(":speech_balloon:  %s\n", tmp)
	}

	return message
}

func GenerateResult(pollID int64, isAdmin bool) string {
	var message string

	var res string
	var user string
	rows, _ := database.Query("select message, author from answer where poll_id = ?", pollID)
	for rows.Next() {
		rows.Scan(&res, &user)
		if isAdmin {
			message += fmt.Sprintf(":speech_balloon:  %s (@%s)\n", res, user)
		} else {
			message += fmt.Sprintf(":speech_balloon:  %s\n", res)
		}
	}

	return message
}

// GetPoll returns the poll id of a specific ts
func GetPoll(ts string) (int64, string) {
	var id int64
	var channel string
	err := database.QueryRow("select id, channel from poll where ts = ?", ts).Scan(&id, &channel)
	if err != nil {
		fmt.Printf("failed to get poll: %v\n", err)
	}
	return id, channel
}

// AddPoll is to insert a new poll in db
func AddPoll(message, channel string) (int64, error) {
	statement, _ := database.Prepare("INSERT INTO poll (message, channel) VALUES (?, ?)")
	res, _ := statement.Exec(message, channel)
	lastID, _ := res.LastInsertId()
	return lastID, nil
}

// AddAnswer adds new answer in database
func AddAnswer(pollID int64, content, author string) {
	statement, _ := database.Prepare("INSERT INTO answer (poll_id, message, author) values (? , ?, ?)")
	statement.Exec(pollID, content, author)
	fmt.Printf("added a new answer into databbase\n")
}

// UpdatePollTS is used to insert th ts after the message has been posted
func UpdatePollTS(id int64, ts string) {
	statement, _ := database.Prepare("update poll set ts=? where id=?")
	statement.Exec(ts, id)
}

// Close can be used to close your database
func Close() {
	database.Close()
}
