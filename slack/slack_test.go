package slack_test

import (
	"os"
	"testing"

	"gitlab-mr-notifier/models"
	"gitlab-mr-notifier/slack"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	require.NoError(t, godotenv.Load("../.env"))

	url := os.Getenv("SLACK_WEBHOOK_URL")
	require.NotEmpty(t, url)

	sl := slack.New()
	err := sl.Notify(url, nil)

	require.NoError(t, err)

	mrs := models.MergeRequests{{
		Title:                         "Mr 1",
		Author:                        "Author 1",
		Description:                   "Description 1",
		UserNotesCount:                1,
		URL:                           "https://gitlab.com/testingapi3/docker-test/-/merge_requests/2",
		HasConflicts:                  false,
		IsBlockingDiscussionsResolved: false,
	}, {
		Title:                         "Mr 2",
		Author:                        "Author 2",
		Description:                   "Description 2",
		UserNotesCount:                123,
		URL:                           "https://gitlab.com/testingapi3/docker-test/-/merge_requests/3",
		HasConflicts:                  true,
		IsBlockingDiscussionsResolved: true,
	}}

	err = sl.Notify(url, mrs)
	require.NoError(t, err)
}

func TestSendErrors(t *testing.T) {
}
