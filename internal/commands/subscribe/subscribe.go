package subscribe

import (
	"fmt"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"github.com/slack-go/slack"
	"log"
)

type Subscriber struct {
	commands.Command
	storage storage.Mapper
}

type args struct {
	id *string
}

func New(s storage.Mapper) *Subscriber {
	return &Subscriber{
		Command: commands.Command{
			Usage: "{project name or group} (ex. 'galaxy/eclipse')",
		},
		storage: s,
	}
}

func (s *Subscriber) Execute(c slack.SlashCommand) slack.Msg {
	args, err := s.parseArgs(c.Text)
	if err != nil {
		return s.Command.Help(err.Error())
	}

	err = s.processSubscription(*args.id, c)
	if err != nil {
		s.Command.Help(err.Error())
	}

	return s.Command.TextMessage("Successfully subscribed")
}

func (s *Subscriber) processSubscription(id string, c slack.SlashCommand) error {
	err := s.storage.Store(id, c.ChannelID, commands.Subscribe)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error save subscription to mapper")
	}
	return nil
}

func (s *Subscriber) parseArgs(input string) (args, error) {
	pieces := s.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) != 1 {
		return res, fmt.Errorf("invalid number of arguments")
	}
	res.id = &pieces[0]
	return res, nil
}
