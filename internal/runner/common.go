package runner

import (
	"log"

	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/gitlabapi"
	lognotifier "gitlab-mr-notifier/internal/logNotifier"
	"gitlab-mr-notifier/internal/models"
	"gitlab-mr-notifier/internal/slack"
)

func check(config config.Config) {
	f := slack.NewSimpleMessageFormatter(config.MessageDescriptionLimit, config.ShortMsgAuthors)
	notifier := slack.New(f)
	if config.Notifier == "log" {
		notifier = lognotifier.New(f)
	}
	gitapi := gitlabapi.New()

	log.Println("start checking")
	var allMRs models.MergeRequests

	for _, id := range config.GitlabGroupIDS {
		mrs, err := gitapi.GetGroupMRList(config.GitlabToken, id)
		if err != nil {
			log.Println("getting mr list error:", err)
		}

		allMRs = append(allMRs, mrs...)
	}

	for _, id := range config.GitlabProjectIDS {
		mrs, err := gitapi.GetProjectMRList(config.GitlabToken, id)
		if err != nil {
			log.Println("getting mr list error:", err)
		}

		allMRs = append(allMRs, mrs...)
	}

	log.Println("found ", len(allMRs), " MRs")
	err := notifier.Notify(config.SlackWebhookURL, allMRs)
	if err != nil {
		log.Println("slack notification error:", err)
	}
}
