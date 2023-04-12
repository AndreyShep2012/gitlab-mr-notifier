package runner

import (
	"log"

	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/gitlabapi"
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
	"gitlab-mr-notifier/internal/slack"

	"github.com/xanzy/go-gitlab"
)

func check(config config.Config) {
	gitapi, err := gitlabapi.New(config.GitlabToken)
	if err != nil {
		log.Println("create gitlabapi error:", err)
		return
	}

	log.Println("start checking MR list")

	mrList, err := gitapi.GetMRList(config.GitlabGroupID)
	if err != nil {
		log.Println("getting mr list error:", err)
		return
	}

	log.Println("found ", len(mrList), " MRs")

	if len(mrList) == 0 {
		return
	}

	var notifyMRs models.MergeRequests

	for _, mr := range mrList {
		if mr.Draft {
			continue
		}

		notifyMR := toModel(mr)
		notifyMR.UnresolvedThreads = getUnresolvedThreads(gitapi, mr.ProjectID, mr.IID)

		notifyMRs = append(notifyMRs, notifyMR)
	}

	sl := slack.New()
	err = sl.Notify(config.SlackWebhookURL, notifyMRs)
	if err != nil {
		log.Println("slack notification error:", err)
	}
}

func toModel(m *gitlab.MergeRequest) models.MergeRequest {
	return models.MergeRequest{
		Title:               m.Title,
		Description:         m.Description,
		Author:              m.Author.Name,
		URL:                 m.WebURL,
		HasConflicts:        m.HasConflicts,
		DetailedMergeStatus: m.DetailedMergeStatus,
		CreatedAt:           *m.CreatedAt,
		UpdatedAt:           *m.UpdatedAt,
	}
}

func getUnresolvedThreads(gitapi interfaces.GitlabApi, projectID, mergeRequestID int) int {
	discussions, err := gitapi.GetMRDiscussions(projectID, mergeRequestID, 0, 100)
	if err != nil {
		log.Printf("getting project[%d] mr[%d] discussions error: %v\n", projectID, mergeRequestID, err)
		return 0
	}

	m := make(map[string]struct{}, len(discussions))

	for _, discussion := range discussions {
		for _, note := range discussion.Notes {
			if note.Resolvable && !note.Resolved {
				m[discussion.ID] = struct{}{}
			}
		}
	}

	return len(m)
}
