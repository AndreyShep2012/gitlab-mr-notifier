package slack

import (
	"errors"
	"fmt"
	"gitlab-mr-notifier/interfaces"
	"gitlab-mr-notifier/models"

	slack_api "github.com/easonlin404/go-slack"
)

type slack struct {
}

func New() interfaces.Notifier {
	return new(slack)
}

func (s slack) Notify(creds interface{}, mrs models.MergeRequests) error {
	webhookURL, ok := creds.(string)
	if !ok {
		return errors.New("wrong credentials format, please use string for webhook url")
	}

	text := "Some MRs are still need to be reviewed:"
	api := slack_api.New().WebhookURL(webhookURL)

	if len(mrs) == 0 {
		return nil
	}

	_, err := api.SendMessage(slack_api.Message{Text: text})
	if err != nil {
		return err
	}

	for _, m := range mrs {
		text = fmt.Sprintf("```Author: %s\nTitle: %s\nURL: %s\nDescription: %s\nHas conflicts: %v\nUserNotesCount: %d\nIsBlockingDiscussionsResolved: %v```",
			m.Author, m.Title, m.URL, m.Description, m.HasConflicts, m.UserNotesCount, m.IsBlockingDiscussionsResolved)
		_, err := api.SendMessage(slack_api.Message{Text: text})
		if err != nil {
			return err
		}
	}

	return nil
}
