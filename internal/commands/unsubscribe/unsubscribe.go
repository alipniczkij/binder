package unsubscribe

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"log"
)

type Unsubscriber struct {
	commands.Command
	storage storage.Mapper
}

type args struct {
	key string
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
	args, err := s.parseArgs(c.Text)
	if err != nil {
		return s.Command.Help(err.Error())
	}
	err = s.storage.Delete(args.key, c.ChannelID, commands.Subscribe)
	if err != nil {
		log.Println(err)
		return s.Command.Help(err.Error())
	}
	return s.Command.TextMessage("Successfully unsubscribed")
}

func (s *Unsubscriber) parseArgs(input string) (args, error) {
	pieces := s.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) != 1 {
		return res, fmt.Errorf("invalid number of arguments")
	}
	res.key = pieces[0]
	return res, nil
}
