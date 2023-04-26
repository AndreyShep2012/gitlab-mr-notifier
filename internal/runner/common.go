package runner

import (
	"log"

	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/gitlabapi"
	"gitlab-mr-notifier/internal/slack"
)

func check(config config.Config) {
	sl := slack.New()
	gitapi := gitlabapi.New()

	log.Println("start checking")
	mrs, err := gitapi.GetMRList(config.GitlabToken, config.GitlabGroupID)
	if err != nil {
		log.Println("getting mr list error:", err)
		return
	}
	log.Println("found ", len(mrs), " MRs")
	err = sl.Notify(config.SlackWebhookURL, mrs)
	if err != nil {
		log.Println("slack notification error:", err)
	}
}
