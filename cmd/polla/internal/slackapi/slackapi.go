package slackapi

import (
	"os"

	"github.com/nlopes/slack"
)

// SlackClient is my slack client
var SlackClient = slack.New(os.Getenv("TEST_POLLA"))

// IsAdmin allows me to check if a user is admin or not.
func IsAdmin(userID string) (bool, error) {
	user, err := SlackClient.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	return !user.IsAdmin, nil
}
