package gitlabapi

import (
	"gitlab-mr-notifier/internal/interfaces"

	"github.com/xanzy/go-gitlab"
)

type gitlabapi struct {
	token  string
	client *gitlab.Client
}

func New(token string) (interfaces.GitlabApi, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &gitlabapi{token: token, client: client}, nil
}

func (ga gitlabapi) GetMRList(groupid int) ([]*gitlab.MergeRequest, error) {
	mrs, _, err := ga.client.MergeRequests.ListGroupMergeRequests(
		groupid,
		&gitlab.ListGroupMergeRequestsOptions{
			State: gitlab.String("opened"),
			Scope: gitlab.String("all"),
			Sort:  gitlab.String("asc"),
		},
	)

	return mrs, err
}

func (ga gitlabapi) GetMRDiscussions(projectID, mergeRequestID, page, perPage int) ([]*gitlab.Discussion, error) {
	discussions, _, err := ga.client.Discussions.ListMergeRequestDiscussions(
		projectID,
		mergeRequestID,
		&gitlab.ListMergeRequestDiscussionsOptions{
			Page:    page,
			PerPage: perPage,
		},
	)

	return discussions, err
}
