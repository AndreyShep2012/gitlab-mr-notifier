package interfaces

import "gitlab-mr-notifier/internal/models"

type Notifier interface {
	Notify(creds interface{}, mrs models.MergeRequests) error
}
