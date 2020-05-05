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

	"github.com/jgengo/Polla/internal/utils"
)

func newPoll(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while Parsing")
		return
	}
	fmt.Printf("%+v\n\n\n", req)

	// text := req.FormValue("text")
	user := req.FormValue("user_id")
	channelID := req.FormValue("channel_id")
	triggerID := req.FormValue("trigger_id")
	responseURL := req.FormValue("response_url")

	isAdmin, err := utils.IsAdmin(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while trying to get user information")
		return
	}
	if !isAdmin {
		utils.ReturnUnauthorized(responseURL)
	}
	utils.NewPollDialog(triggerID)
	go func() {
		time.Sleep(time.Second * 10)
		utils.SendPoll(channelID)
		time.Sleep(time.Second * 5)
		utils.UpdateLastPoll()
	}()

}

func interactivity(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while Parsing")
		return
	}
	fmt.Printf("new interactivity: %+v\n", req)
}

func root(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error while Parsing")
		return
	}
	fmt.Printf("new call: %+v\n", req)
}

func main() {
	srv := &http.Server{
		Handler:      nil,
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/interactivity", interactivity)
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
