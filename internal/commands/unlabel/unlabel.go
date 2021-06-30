package unlabel

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"log"
)

type Unlabeler struct {
	commands.Command
	storage storage.Mapper
}

type args struct {
	labels []string
}

func New(s storage.Mapper) *Unlabeler {
	return &Unlabeler{
		Command: commands.Command{
			Usage: "{label that you want to exclude} (ex. 'backend frontend')",
		},
		storage: s,
	}
}

func (l *Unlabeler) Execute(c slack.SlashCommand) slack.Msg {
	args, err := l.parseArgs(c.Text)
	if err != nil {
		return l.Command.Help(err.Error())
	}
	err = l.processLabel(args, c)
	if err != nil {
		l.Command.Help(err.Error())
	}

	return l.Command.TextMessage("Labels excluded")
}

func (l *Unlabeler) processLabel(args args, c slack.SlashCommand) error {
	for _, label := range args.labels {
		err := l.storage.Store(c.ChannelID, label, commands.Unlabel)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error save label to mapper")
		}
	}
	return nil
}

func (l *Unlabeler) parseArgs(input string) (args, error) {
	pieces := l.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) < 1 {
		return res, fmt.Errorf("labels not provided")
	}
	res.labels = pieces
	return res, nil
}
