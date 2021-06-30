package label

import (
	"fmt"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"github.com/slack-go/slack"
	"log"
)

type Labeler struct {
	commands.Command
	storage storage.Mapper
}

type args struct {
	labels []string
}

func New(s storage.Mapper) *Labeler {
	return &Labeler{
		Command: commands.Command{
			Usage: "{label that you want to include} (ex. 'backend frontend')",
		},
		storage: s,
	}
}

func (l *Labeler) Execute(c slack.SlashCommand) slack.Msg {
	args, err := l.parseArgs(c.Text)
	if err != nil {
		return l.Command.Help(err.Error())
	}
	err = l.processLabel(args, c)
	if err != nil {
		l.Command.Help(err.Error())
	}

	return l.Command.TextMessage("Labels added")
}

func (l *Labeler) processLabel(args args, c slack.SlashCommand) error {
	for _, label := range args.labels {
		err := l.storage.Store(c.ChannelID, label, commands.Label)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error save label to mapper")
		}
	}
	return nil
}

func (l *Labeler) parseArgs(input string) (args, error) {
	pieces := l.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) < 1 {
		return res, fmt.Errorf("labels not provided")
	}
	res.labels = pieces
	return res, nil
}
