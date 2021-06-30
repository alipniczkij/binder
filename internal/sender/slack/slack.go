package slack

import (
	"fmt"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/gitlab"
	"github.com/alipniczkij/binder/internal/storage"
	"github.com/alipniczkij/binder/pkg/config"
	"github.com/alipniczkij/binder/pkg/models"
	"github.com/slack-go/slack"
	"log"
	"net/url"
	"strings"
)

type Sender struct {
	slack  *slack.Client
	mapper storage.Mapper
}

func New(cfg config.Slack, mapper storage.Mapper) *Sender {
	slackClient := slack.New(cfg.Token)
	return &Sender{slack: slackClient, mapper: mapper}
}

func (s *Sender) Send(event *models.GitlabEvent) error {
	data := s.processEvent(event)
	if data == nil {
		return nil
	}
	for channel := range data.channels {
		if !s.check(event, channel) {
			continue
		}
		channelID, _, err := s.slack.PostMessage(
			channel,
			slack.MsgOptionText(data.text, false),
			slack.MsgOptionAsUser(true),
		)
		if err != nil {
			log.Printf("Can't send message to %s: %s\n", channelID, err.Error())
		}
	}
	return nil
}

func (s *Sender) processEvent(event *models.GitlabEvent) *Data {
	if event.ObjectAttributes.Action != "open" {
		return nil
	}
	u, err := url.Parse(event.ObjectAttributes.URL)
	if err != nil {
		log.Printf("can't get project ID from url %s", event.ObjectAttributes.URL)
		return nil
	}
	projectID, err := gitlab.GetProjectIDFromURL(u)
	if err != nil {
		log.Printf("can't get project ID from url %s", event.ObjectAttributes.URL)
		return nil
	}

	channels := s.getChannels(projectID)

	return &Data{
		text:     fmt.Sprintf("%s has opened new <%s|issue>", event.User.Name, event.ObjectAttributes.URL),
		channels: channels,
	}
}

func (s *Sender) getChannels(project string) map[string]bool {
	channels := make(map[string]bool)
	pieces := strings.Split(project, "/")
	for i := 1; i < len(pieces)+1; i++ {
		subs, found := s.mapper.Get(strings.Join(pieces[:i], "/"), commands.Subscribe)
		if !found {
			continue
		}
		for _, sub := range subs {
			channels[sub] = true
		}
	}
	return channels
}

func (s *Sender) check(event *models.GitlabEvent, channel string) bool {
	res, found := s.mapper.Get(channel, commands.Unlabel)
	if found {
		for _, label := range event.Labels {
			if contains(label.Title, res) {
				return false
			}
		}
	}
	res, found = s.mapper.Get(channel, commands.Label)
	if found {
		for _, label := range event.Labels {
			if contains(label.Title, res) {
				return true
			}
		}
		return false
	}
	return true
}

func contains(l string, ss []string) bool {
	for _, s := range ss {
		if l == s {
			return true
		}
	}
	return false
}
