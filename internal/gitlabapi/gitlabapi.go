package gitlabapi

import (
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"

	"github.com/xanzy/go-gitlab"
)

type gitlabapi struct {
}

func New() interfaces.GitlabApi {
	return new(gitlabapi)
}

func (ga gitlabapi) GetMRList(token string, groupid int) (models.MergeRequests, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	mrs, _, err := client.MergeRequests.ListGroupMergeRequests(groupid, &gitlab.ListGroupMergeRequestsOptions{
		State: gitlab.String("opened"),
		Scope: gitlab.String("all"),
	})

	return toModels(mrs), err
}

func toModels(mrs []*gitlab.MergeRequest) models.MergeRequests {
	var res models.MergeRequests
	for _, m := range mrs {
		if m.Draft {
			continue
		}
		res = append(res, toModel(m))
	}
	return res
}

func toModel(m *gitlab.MergeRequest) models.MergeRequest {
	return models.MergeRequest{
		Title:                         m.Title,
		Description:                   m.Description,
		Author:                        m.Author.Name,
		UserNotesCount:                m.UserNotesCount,
		URL:                           m.WebURL,
		HasConflicts:                  m.HasConflicts,
		IsBlockingDiscussionsResolved: m.BlockingDiscussionsResolved,
	}
}
