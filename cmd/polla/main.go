package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jgengo/Polla/internal/db"
	"github.com/jgengo/Polla/internal/utils"
	"github.com/slack-go/slack"
)

func command(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while Parsing")
		return
	}

	user := req.FormValue("user_id")
	triggerID := req.FormValue("trigger_id")

	isAdmin, err := utils.IsAdmin(user)
	if err != nil {
		log.Printf("[error] failed to check if user is admin: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while trying to get user information")
		return
	}
	if !isAdmin {
		fmt.Fprintf(w, "Sorry, you are not authorized to use this command")
		return
	}
	utils.NewPollDialog(triggerID)
}

func interactivity(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var message slack.InteractionCallback

	buf, _ := ioutil.ReadAll(req.Body)
	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[error] failed to unespace request body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[error] failed to decode json message from slack: %s\n", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if message.Type == "dialog_submission" {
		channel := message.Channel.GroupConversation.Conversation.ID
		content := message.Submission["content"]
		if message.CallbackID == "new_poll" {
			utils.SendPoll(channel, content)
		}
		if len(message.CallbackID) > 10 && message.CallbackID[:10] == "new_answer" {
			ts := message.CallbackID[11:]
			userID := message.User.Name
			utils.SendAnswer(ts, content, userID)
		}
	}

	if message.Type == "block_actions" {
		actions := message.ActionCallback.BlockActions

		actionID := actions[0].ActionID
		if actionID == "submit" {
			utils.NewAnswerDialog(message.TriggerID, message.Message.Timestamp)
		}
		if actionID == "result" {
			utils.ShowResults(message.User.ID, message.Message.Timestamp)
		}
	}
}

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("failed to init the database")
	}

	srv := &http.Server{
		Handler:      nil,
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/interactivity", interactivity)
	http.HandleFunc("/command", command)

	go func() {
		log.Println("Starting Server")
		if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	db.Close()

	log.Println("Shutting down")
	os.Exit(0)

}
