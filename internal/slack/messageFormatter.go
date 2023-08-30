package slack

import (
	"fmt"
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
)

type messageFormatter struct {
	limit int
}

const SimpleMessageFormatterNoLimit = 0

func NewSimpleMessageFormatter(limit int) interfaces.MessageFormatter {
	return &messageFormatter{limit: limit}
}

func (m messageFormatter) GetBody(mr models.MergeRequest) string {
	if m.limit > SimpleMessageFormatterNoLimit {
		mr.Description = crop(mr.Description, m.limit)
	}

	const defaultTimeLayout = "2006-01-02 15:04:05"
	return fmt.Sprintf("```Author: %s\nTitle: %s\nURL: %s\nDescription: %s\n\nHasConflicts: %v\nDetailedMergeStatus: %s\nUnresolvedThreads: %d\nCreatedAt: %s\nUpdatedAt: %s```",
		mr.Author, mr.Title, mr.URL, mr.Description, mr.HasConflicts, mr.DetailedMergeStatus, mr.UnresolvedThreads, mr.CreatedAt.Format(defaultTimeLayout), mr.UpdatedAt.Format(defaultTimeLayout))
}

func (m messageFormatter) GetIntroText(mrsCount int) string {
	if mrsCount < 0 {
		return fmt.Sprintf("Wrong number of MRS: %d !!!", mrsCount)
	} else if mrsCount == 0 {
		return "Hooray. No MRs to review!"
	} else if mrsCount == 1 {
		return "1 MR is still need to be reviewed:"
	}

	return fmt.Sprintf("%d MRs are still need to be reviewed:", mrsCount)
}

func crop(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}

	return s
}
