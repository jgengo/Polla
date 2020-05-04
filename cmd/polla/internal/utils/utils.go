package utils

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

// SlackClient is my slack client
var SlackClient = slack.New(os.Getenv("TEST_POLLA"))

// IsAdmin allows me to check if a user is admin or not.
func IsAdmin(userID string) (bool, error) {
	user, err := SlackClient.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}

// ReturnUnauthorized webhook
func ReturnUnauthorized(url string) {
	resp := &slack.WebhookMessage{
		Text: "Sorry, you are not authorized to use this command",
	}
	slack.PostWebhook(url, resp)
}

func StartDialog(triggerID, userID string) {

	dg := slack.NewTextInput("test", "test2", "text test")
	var ddg []slack.DialogElement

	ddg = append(ddg, dg)

	dialog := slack.Dialog{
		TriggerID:  triggerID,
		CallbackID: "abc",
		Title:      "Add a new Poll",
		Elements:   ddg,
	}
	if err := SlackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
