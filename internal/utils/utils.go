package utils

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var lastChannelID string
var lastTS string

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

func dialogAddAnswer() {

}

// SlackClient is my slack client
var SlackClient = slack.New(os.Getenv("TEST_POLLA"))

// IsAdmin returns if a specified user is admin or not
func IsAdmin(userID string) (bool, error) {
	user, err := SlackClient.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}

// ReturnUnauthorized returns unauthorize via quick-response webhook
func ReturnUnauthorized(url string) {
	resp := &slack.WebhookMessage{
		Text: "Sorry, you are not authorized to use this command",
	}
	slack.PostWebhook(url, resp)
}

// NewPollDialog will send the dialog to add a new poll
func NewPollDialog(triggerID string) {
	dialog := dialogNewPoll(triggerID)
	if err := SlackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}

func SendPoll(channelID string) {
	headerText := slack.NewTextBlockObject("mrkdwn", "Question?", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	newBtnTxt := slack.NewTextBlockObject("plain_text", "Submit Response", false, false)
	newBtn := slack.NewButtonBlockElement("", "click_me_123", newBtnTxt)
	actionBlock := slack.NewActionBlock("", newBtn)

	_, ts, err := SlackClient.PostMessage(channelID, slack.MsgOptionText("New Poll started!", false), slack.MsgOptionBlocks(headerSection, actionBlock))
	if err != nil {
		fmt.Printf("error pushing: %+v\n", err)
	}

	lastChannelID = channelID
	lastTS = ts
}

func UpdateLastPoll() {

	headerText := slack.NewTextBlockObject("mrkdwn", "Question?\n\n:speech_balloon: When do we eat?", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	newBtnTxt := slack.NewTextBlockObject("plain_text", "Submit Response", false, false)
	newBtn := slack.NewButtonBlockElement("", "click_me_123", newBtnTxt)
	actionBlock := slack.NewActionBlock("", newBtn)

	_, _, _, err := SlackClient.UpdateMessage(lastChannelID, lastTS, slack.MsgOptionText("New Poll started!", false), slack.MsgOptionBlocks(headerSection, actionBlock))
	if err != nil {
		fmt.Printf("error updating: %+v\n\n", err)
	}
}

/*
{
	"blocks": [
		{
			"type": "section",
			"text": {"type": "mrkdwn", "text": "Question\n\n" }
		},
		{
			"type": "actions",
			"elements": [
				{
					"type": "button",
					"text": {
						"type": "plain_text",
						"text": "Submit Response",
						"emoji": true
					},
					"value": "click_me_123"
				}
			]
		}
	]
}
*/
