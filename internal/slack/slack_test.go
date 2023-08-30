package slack_test

import (
	"os"
	"testing"
	"time"

	"gitlab-mr-notifier/internal/models"
	"gitlab-mr-notifier/internal/slack"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	url := os.Getenv("SLACK_WEBHOOK_URL")
	require.NotEmpty(t, url)

	sl := slack.New(slack.NewSimpleMessageFormatter(slack.SimpleMessageFormatterNoLimit))
	err := sl.Notify(url, nil)

	require.NoError(t, err)

	mrs := models.MergeRequests{
		{
			Title:               "Mr 1",
			Author:              "Author 1",
			Description:         "Description 1",
			URL:                 "https://gitlab.com/testingapi3/docker-test/-/merge_requests/2",
			HasConflicts:        false,
			UnresolvedThreads:   0,
			DetailedMergeStatus: "not_approved",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			Title:               "Mr 2",
			Author:              "Author 2",
			Description:         "Description 2",
			URL:                 "https://gitlab.com/testingapi3/docker-test/-/merge_requests/3",
			HasConflicts:        true,
			UnresolvedThreads:   1,
			DetailedMergeStatus: "not_approved",
			CreatedAt:           time.Now().AddDate(0, 0, -2),
			UpdatedAt:           time.Now().AddDate(0, 0, -1),
		},
	}

	err = sl.Notify(url, mrs)
	require.NoError(t, err)
}

func TestSimpleFormatterBody(t *testing.T) {
	f := slack.NewSimpleMessageFormatter(-100)
	require.Contains(t, f.GetBody(models.MergeRequest{Description: "desc1"}), "desc1")

	f = slack.NewSimpleMessageFormatter(slack.SimpleMessageFormatterNoLimit)
	require.Contains(t, f.GetBody(models.MergeRequest{Description: "desc1"}), "desc1")

	f = slack.NewSimpleMessageFormatter(5)
	body := f.GetBody(models.MergeRequest{Description: "desc123"})
	require.Contains(t, body, "desc1")
	require.NotContains(t, body, "desc12")

	f = slack.NewSimpleMessageFormatter(5000)
	require.Contains(t, f.GetBody(models.MergeRequest{Description: "desc1"}), "desc1")
}

func TestSimpleFormatterIntro(t *testing.T) {
	f := slack.NewSimpleMessageFormatter(slack.SimpleMessageFormatterNoLimit)
	require.Equal(t, "Hooray. No MRs to review!", f.GetIntroText(0))
	require.Equal(t, "Wrong number of MRS: -10 !!!", f.GetIntroText(-10))
	require.Equal(t, "1 MR is still need to be reviewed:", f.GetIntroText(1))
	require.Equal(t, "10 MRs are still need to be reviewed:", f.GetIntroText(10))
}
