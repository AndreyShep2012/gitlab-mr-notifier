package interfaces

import "gitlab-mr-notifier/internal/models"

type MessageFormatter interface {
	GetBody(models.MergeRequest) string
	GetIntroText(mrs int) string
}
