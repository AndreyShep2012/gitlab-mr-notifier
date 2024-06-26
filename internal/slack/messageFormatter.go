package slack

import (
	"fmt"
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
	"slices"
	"strings"
)

type messageFormatter struct {
	limit               int
	shortMessageAuthors []string
}

const SimpleMessageFormatterNoLimit = 0
const defaultTimeLayout = "2006-01-02 15:04:05"

func NewSimpleMessageFormatter(limit int, shortMessageAuthors string) interfaces.MessageFormatter {
	return &messageFormatter{limit: limit, shortMessageAuthors: strings.Split(shortMessageAuthors, ",")}
}

func (m messageFormatter) GetBody(mr models.MergeRequest) string {
	if m.limit > SimpleMessageFormatterNoLimit {
		mr.Description = crop(mr.Description, m.limit)
	}

	if slices.Contains(m.shortMessageAuthors, mr.Author) {
		return getBodyNoDescription(mr)
	}

	return getBodyWithDescription(mr)
}

func getBodyWithDescription(mr models.MergeRequest) string {
	return fmt.Sprintf("```Author: %s\nTitle: %s\nURL: %s\nDescription: %s\n\nHasConflicts: %v\nDetailedMergeStatus: %s\nUnresolvedThreads: %d\nBranch: %s\nChangesCount: %s\nCreatedAt: %s\nUpdatedAt: %s```",
		mr.Author, mr.Title, mr.URL, mr.Description, mr.HasConflicts, mr.DetailedMergeStatus, mr.UnresolvedThreads, mr.Branch, mr.ChangesCount, mr.CreatedAt.Format(defaultTimeLayout), mr.UpdatedAt.Format(defaultTimeLayout))
}

func getBodyNoDescription(mr models.MergeRequest) string {
	return fmt.Sprintf("```Author: %s\nTitle: %s\nURL: %s\n\nUnresolvedThreads: %d\nBranch: %s\nChangesCount: %s\nCreatedAt: %s```",
		mr.Author, mr.Title, mr.URL, mr.UnresolvedThreads, mr.Branch, mr.ChangesCount, mr.CreatedAt.Format(defaultTimeLayout))
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

func (m messageFormatter) GetPipelineFailedIntroText(mrs int) string {
	return "MRs with failed pipeline"
}

func (m messageFormatter) GetPipelineFailedGetBody(mrs []models.MergeRequest) string {
	var sb strings.Builder
	sb.WriteString("```\n")

	for _, mr := range mrs {
		sb.WriteString(fmt.Sprintf("Author: %s\n", mr.Author))
		sb.WriteString(fmt.Sprintf("Title: %s\n", mr.Title))
		sb.WriteString(fmt.Sprintf("URL: %s\n", mr.URL))
		sb.WriteString(fmt.Sprintf("Pipeline: %s\n\n", mr.PipelineInfo.URL))
	}
	sb.WriteString("```")
	return sb.String()
}

func crop(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}

	return s
}
