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

func (ga gitlabapi) GetGroupMRList(token string, groupid int) (models.MergeRequests, error) {
	return ga.getMRList(token, groupid, getGroupMRListByPage)
}

func (ga gitlabapi) GetProjectMRList(token string, projectid int) (models.MergeRequests, error) {
	return ga.getMRList(token, projectid, getProjectMRListByPage)
}

func (ga gitlabapi) getMRList(token string, id int, f func(*gitlab.Client, int, int) ([]*gitlab.MergeRequest, *gitlab.Response, error)) (models.MergeRequests, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	var res models.MergeRequests
	page := 1
	for {
		mrs, response, err := f(client, id, page)
		if err != nil {
			return res, err
		}

		for _, mr := range mrs {
			fullInfo := getMergeRequest(client, mr.ProjectID, mr.IID)
			if fullInfo != nil {
				mr = fullInfo
			}
			notifyMR := toModel(mr)
			notifyMR.UnresolvedThreads = getUnresolvedThreads(client, mr.ProjectID, mr.IID)

			res = append(res, notifyMR)
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		page++
	}

	return res, err
}

func getGroupMRListByPage(client *gitlab.Client, groupid, page int) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return client.MergeRequests.ListGroupMergeRequests(groupid, &gitlab.ListGroupMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{Page: page},
		State:       gitlab.String("opened"),
		Sort:        gitlab.String("asc"),
		WIP:         gitlab.String("no"),
	})
}

func getProjectMRListByPage(client *gitlab.Client, projectId, page int) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{Page: page},
		State:       gitlab.String("opened"),
		Sort:        gitlab.String("asc"),
		WIP:         gitlab.String("no"),
	})
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

func getMergeRequest(client *gitlab.Client, projectID, mergeRequestID int) *gitlab.MergeRequest {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mergeRequestID, nil)
	if err != nil {
		log.Printf("getting project[%d] mr[%d] full info error: %v\n", projectID, mergeRequestID, err)
		return nil
	}

	return mr
}

func toModel(m *gitlab.MergeRequest) models.MergeRequest {
	mr := models.MergeRequest{
		Title:               m.Title,
		Description:         m.Description,
		Author:              m.Author.Name,
		URL:                 m.WebURL,
		HasConflicts:        m.HasConflicts,
		DetailedMergeStatus: m.DetailedMergeStatus,
		CreatedAt:           *m.CreatedAt,
		UpdatedAt:           *m.UpdatedAt,
		ChangesCount:        m.ChangesCount,
		Branch:              m.SourceBranch,
	}

	if m.Pipeline != nil {
		mr.PipelineInfo.IsFailed = m.Pipeline.Status == "failed"
		mr.PipelineInfo.URL = m.Pipeline.WebURL
	}

	return mr
}
