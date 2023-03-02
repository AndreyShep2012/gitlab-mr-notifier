package gitapi

import (
	"gitlab-mr-notifier/interfaces"
	"gitlab-mr-notifier/models"

	"github.com/xanzy/go-gitlab"
)

type gitapi struct {
}

func New() interfaces.GitApi {
	return new(gitapi)
}

func (ga gitapi) GetMRList(token string, groupid int) (models.MergeRequests, error) {
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
