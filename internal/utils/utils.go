package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jgengo/Polla/internal/db"
	"github.com/slack-go/slack"
)

// var lastChannelID string
// var lastTS string

func dialogNewPoll(triggerID string) slack.Dialog {
	dg := slack.NewTextInput("content", "Question", "")
	dg.MaxLength = 150
	dg.Placeholder = "Write something"
	var ddg []slack.DialogElement

	ddg = append(ddg, dg)

	dialog := slack.Dialog{
		TriggerID:   triggerID,
		CallbackID:  "new_poll",
		Title:       "Add a new Poll",
		SubmitLabel: "Create",
		Elements:    ddg,
	}

	return dialog
}

func dialogNewAnser(triggerID, messageTS string) slack.Dialog {
	dg := slack.NewTextInput("content", "Answer", "")
	dg.MaxLength = 150
	dg.Placeholder = "Write something"
	var ddg []slack.DialogElement

	ddg = append(ddg, dg)

	dialog := slack.Dialog{
		TriggerID:   triggerID,
		CallbackID:  "new_answer:" + messageTS,
		Title:       "Add a Response",
		SubmitLabel: "Submit",
		Elements:    ddg,
	}

	return dialog
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

// NewPollDialog will send the dialog to add a new poll
func NewPollDialog(triggerID string) {
	dialog := dialogNewPoll(triggerID)
	if err := SlackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}

func NewAnswerDialog(triggerID, messageTS string) {
	dialog := dialogNewAnser(triggerID, messageTS)
	if err := SlackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
func SendPoll(channelID, question string) {
	dbID, _ := db.AddPoll(question, channelID)
	dbIDStr := strconv.FormatInt(dbID, 10)

	headerText := slack.NewTextBlockObject("mrkdwn", question, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	newBtnTxt := slack.NewTextBlockObject("plain_text", "Submit Response", false, false)
	newBtn := slack.NewButtonBlockElement("submit", dbIDStr, newBtnTxt)
	actionBlock := slack.NewActionBlock("", newBtn)

	_, ts, err := SlackClient.PostMessage(channelID, slack.MsgOptionText("New Poll started!", false), slack.MsgOptionBlocks(headerSection, actionBlock))
	if err != nil {
		fmt.Printf("error pushing: %+v\n", err)
	}

	fmt.Printf("== message TS: %s\n", ts)

	db.UpdatePollTS(dbID, ts)
}

func SendAnswer(ts, content, userID string) {
	pollID, channelID := db.GetPoll(ts)
	dbIDStr := strconv.FormatInt(pollID, 10)

	db.AddAnswer(pollID, content, userID)
	newTxt := db.GenerateText(pollID)

	headerText := slack.NewTextBlockObject("mrkdwn", newTxt, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	newBtnTxt := slack.NewTextBlockObject("plain_text", "Submit Response", false, false)
	newBtn := slack.NewButtonBlockElement("submit", dbIDStr, newBtnTxt)
	actionBlock := slack.NewActionBlock("", newBtn)
	_, _, _, err := SlackClient.UpdateMessage(channelID, ts, slack.MsgOptionText("New Poll started!", false), slack.MsgOptionBlocks(headerSection, actionBlock))
	if err != nil {
		fmt.Printf("error updating: %+v\n\n", err)
	}
}

// func UpdateLastPoll() {
// 	headerText := slack.NewTextBlockObject("mrkdwn", "Question?\n\n:speech_balloon: When do we eat?", false, false)
// 	headerSection := slack.NewSectionBlock(headerText, nil, nil)

// 	newBtnTxt := slack.NewTextBlockObject("plain_text", "Submit Response", false, false)
// 	newBtn := slack.NewButtonBlockElement("", "click_me_123", newBtnTxt)
// 	actionBlock := slack.NewActionBlock("", newBtn)

// }

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
