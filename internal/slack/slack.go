package slack

import (
	"errors"
	"fmt"

	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"

	slack_api "github.com/easonlin404/go-slack"
)

type slack struct {
}

func New() interfaces.Notifier {
	return new(slack)
}

func (s slack) Notify(creds interface{}, mrs models.MergeRequests) error {
	var text string

	l := len(mrs)

	switch l {
	case 0:
		return nil
	case 1:
		text = "1 MR is still need to be reviewed:"
	default:
		text = fmt.Sprintf("%d MRs are still need to be reviewed:", l)
	}

	webhookURL, ok := creds.(string)
	if !ok {
		return errors.New("wrong credentials format, please use string for webhook url")
	}

	api := slack_api.New().WebhookURL(webhookURL)

	_, err := api.SendMessage(slack_api.Message{Text: text})
	if err != nil {
		return err
	}

	const defaultTimeLayout = "2006-01-02 15:04:05"

	for _, m := range mrs {
		text = fmt.Sprintf("```Author: %s\nTitle: %s\nURL: %s\nDescription: %s\n\nHasConflicts: %v\nDetailedMergeStatus: %s\nUnresolvedThreads: %d\nCreatedAt: %s\nUpdatedAt: %s```",
			m.Author, m.Title, m.URL, m.Description, m.HasConflicts, m.DetailedMergeStatus, m.UnresolvedThreads, m.CreatedAt.Format(defaultTimeLayout), m.UpdatedAt.Format(defaultTimeLayout))

		_, err := api.SendMessage(slack_api.Message{Text: text})
		if err != nil {
			return err
		}
	}

	return nil
}
