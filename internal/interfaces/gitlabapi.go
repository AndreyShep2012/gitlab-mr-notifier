package interfaces

import "gitlab-mr-notifier/internal/models"

type GitlabApi interface {
	GetProjectMRList(token string, projectId int) (models.MergeRequests, error)
	GetGroupMRList(token string, groupid int) (models.MergeRequests, error)
}
