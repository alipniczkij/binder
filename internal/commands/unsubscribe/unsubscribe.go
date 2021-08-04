package unsubscribe

import (
	"fmt"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"github.com/alipniczkij/binder/pkg/multi_error"
	"github.com/slack-go/slack"
	"log"
)

type Unsubscriber struct {
	commands.Command
	storage storage.Mapper
}

type args struct {
	projects []string
}

func New(s storage.Mapper) *Unsubscriber {
	return &Unsubscriber{
		Command: commands.Command{
			Usage: "{project name or group} (ex. 'galaxy/eclipse')",
		},
		storage: s,
	}
}

func (s *Unsubscriber) Execute(c slack.SlashCommand) slack.Msg {
	errors := multi_error.New()
	args, err := s.parseArgs(c.Text)
	if err != nil {
		return s.Command.Help(err.Error())
	}
	for _, p := range args.projects {
		err = s.storage.Delete(p, c.ChannelID, commands.Subscribe)
		if err != nil {
			log.Println(err)
			errors.Append(err)
		}
	}
	if !errors.IsEmpty() {
		return s.Command.Help(errors.Error())
	}

	return s.Command.TextMessage("Successfully unsubscribed")
}

func (s *Unsubscriber) parseArgs(input string) (args, error) {
	pieces := s.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) < 1 {
		return res, fmt.Errorf("invalid number of arguments")
	}
	res.projects = pieces
	return res, nil
}
