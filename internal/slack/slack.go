package slack

import (
	"errors"

	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"

	slack_api "github.com/easonlin404/go-slack"
)

type slack struct {
	formatter interfaces.MessageFormatter
}

func New(formatter interfaces.MessageFormatter) interfaces.Notifier {
	return &slack{formatter: formatter}
}

func (s slack) Notify(creds interface{}, mrs models.MergeRequests) error {
	webhookURL, ok := creds.(string)
	if !ok {
		return errors.New("wrong credentials format, please use string for webhook url")
	}

	api := slack_api.New().WebhookURL(webhookURL)
	_, err := api.SendMessage(slack_api.Message{Text: s.formatter.GetIntroText(len(mrs))})
	if err != nil {
		return err
	}

	for _, m := range mrs {
		_, err := api.SendMessage(slack_api.Message{Text: s.formatter.GetBody(m)})
		if err != nil {
			return err
		}
	}

	return nil
}
