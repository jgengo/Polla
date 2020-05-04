package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jgengo/Polla/cmd/polla/internal/slackapi"
	"github.com/nlopes/slack"
)

func newPoll(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while Parsing")
	}
	fmt.Printf("%+v\n\n\n", req)

	text := req.FormValue("text")
	user := req.FormValue("user_id")
	responseURL := req.FormValue("response_url")

	fmt.Println(user + " typed /polla " + text)
	fmt.Println(responseURL)

	isAdmin, err := slackapi.IsAdmin(user)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	if !isAdmin {
		resp := &slack.WebhookMessage{
			Text: "Sorry, you are not authorized to use this command",
		}
		slack.PostWebhook(responseURL, resp)
	}

}

func root(w http.ResponseWriter, req *http.Request) {

}

func main() {
	srv := &http.Server{
		Handler:      nil,
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/command", newPoll)
	http.HandleFunc("/", root)

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

	log.Println("Shutting down")
	os.Exit(0)

}
