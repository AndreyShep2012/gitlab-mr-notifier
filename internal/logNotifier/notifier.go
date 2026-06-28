package lognotifier

import (
	"fmt"
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
	"log"
	"os"
	"strings"
)

type notifier struct {
	formatter interfaces.MessageFormatter
}

func New(formatter interfaces.MessageFormatter) interfaces.Notifier {
	return &notifier{formatter: formatter}
}

func (n notifier) Notify(creds any, mrs models.MergeRequests) error {
	var sb strings.Builder
	for _, m := range mrs {
		sb.WriteString(n.formatter.GetBody(m))
		sb.WriteString("\n\n")
	}

	f, err := os.CreateTemp("", "mr-notifier-*.txt")
	if err != nil {
		return fmt.Errorf("logNotifier: create temp file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("logNotifier: write temp file: %w", err)
	}

	log.Println("logNotifier output:", f.Name())
	return nil
}
