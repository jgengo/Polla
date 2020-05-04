package utils

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func dialogNewPoll(triggerID string) slack.Dialog {
	dg := slack.NewTextInput("question", "Question", "")
	dg.MaxLength = 150
	dg.Placeholder = "Write something"
	var ddg []slack.DialogElement

	ddg = append(ddg, dg)

	dialog := slack.Dialog{
		TriggerID:      triggerID,
		CallbackID:     "abc",
		Title:          "Add a new Poll",
		SubmitLabel:    "Create",
		NotifyOnCancel: true,
		Elements:       ddg,
	}

	return dialog
}

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

func NewPollDialog(triggerID string) {
	dialog := dialogNewPoll(triggerID)
	if err := SlackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
