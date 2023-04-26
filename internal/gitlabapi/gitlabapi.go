package gitlabapi

import (
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"
	"log"

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
		Sort:  gitlab.String("asc"),
	})

	res := make(models.MergeRequests, 0, len(mrs))
	for _, mr := range mrs {
		if mr.Draft {
			continue
		}

		notifyMR := toModel(mr)
		notifyMR.UnresolvedThreads = getUnresolvedThreads(client, mr.ProjectID, mr.IID)

		res = append(res, notifyMR)
	}

	return res, err
}

func getUnresolvedThreads(client *gitlab.Client, projectID, mergeRequestID int) int {
	discussions, _, err := client.Discussions.ListMergeRequestDiscussions(
		projectID,
		mergeRequestID,
		&gitlab.ListMergeRequestDiscussionsOptions{
			Page:    0,
			PerPage: 100,
		},
	)
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
