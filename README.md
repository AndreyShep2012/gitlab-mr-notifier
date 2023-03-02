# gitlab-mr-notifier

Service which gets opened MRs from gitlab group periodically and post them into the slack channel. It's used to remind developers about not reviewed MRs

<img width="1008" alt="andreyshep2012 - Test - Slack 2023-03-02 12-19-39" src="https://user-images.githubusercontent.com/30069672/222518158-d605712a-07b3-456f-b4c0-7bb5ee46170e.png">

### Setup Gitlab

Gitlab [read_api](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html#create-a-project-access-token) readonly token is required to receive informations about opened merge requests

Also to check perticullar group it's necessary to setup [group_id](https://docs.gitlab.com/ee/user/group/) in service's settings


### Setup Slack

[Slack webhook](https://api.slack.com/messaging/webhooks) is used to send notifications to the particular channel, so it should be configured before and passed to the service's settings

### Setup cron

Golang [cron](github.com/go-co-op/gocron) library is used under the hood. So if some short period of time is necessary for the cron job, only `CRON_PERIOD` variable can be used 

```
CRON_PERIOD=10s
or
CRON_PERIOD=10m
```

For the long term, like day or week, variable `CRON_TIME` is necessary to setup concrete time

```
CRON_PERIOD=1d
CRON_TIME=10:30
```

### Run locally

Create create text file with environment variables (.env format)

```
GITLAB_TOKEN=`your-token`
GITLAB_GROUP_ID=`your-id`
SLACK_WEBHOOK_URL=`your webhook url`
SLACK_USER=Gitlab MR notifier
CRON_PERIOD=1d
CRON_TIME=10:30
```

Put file path to `CONFIG_PATH` env variable, build app and start service

```
$ make build
$ export CONFIG_PATH="./.env.sample"&&./gitlab-mr-notifier
```


### Testing

Create `.env` file in root with test values

```
GITLAB_TOKEN=`your-token`
GITLAB_GROUP_ID=`your-id`
SLACK_WEBHOOK_URL=`your webhook url`
SLACK_USER=Gitlab MR notifier
CRON_PERIOD=1d
CRON_TIME=10:30
```

Run tests

```
$ make test
```

### Docker

**_WARNING:_** Current approach uses environment variables, maybe it would be better to use Docker secrets instead

Build image:

```
$ docker build -t gitlab-mr-notifier .
```

Create config env file and use your file name to start container, for example:

```
$ docker run -it --rm --env-file=.env gitlab-mr-notifier
```
