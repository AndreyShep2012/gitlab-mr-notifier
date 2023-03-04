package interfaces

import "gitlab-mr-notifier/internal/models"

type GitApi interface {
	GetMRList(token string, groupid int) (models.MergeRequests, error)
}
