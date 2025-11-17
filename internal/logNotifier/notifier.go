package lognotifier

import (
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
	"log"
)

type notifier struct {
	formatter interfaces.MessageFormatter
}

func New(formatter interfaces.MessageFormatter) interfaces.Notifier {
	return &notifier{formatter: formatter}
}

func (n notifier) Notify(creds any, mrs models.MergeRequests) error {
	for _, m := range mrs {
		log.Println(n.formatter.GetBody(m))
	}
	return nil
}
