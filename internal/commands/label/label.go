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
	delete bool
}

func New(s storage.Mapper) *Labeler {
	return &Labeler{
		Command: commands.Command{
			Usage: "{labels that you want to include} (ex. 'backend frontend')\n -delete {labels}",
		},
		storage: s,
	}
}

func (l *Labeler) Execute(c slack.SlashCommand) slack.Msg {
	args, err := l.parseArgs(c.Text)
	if err != nil {
		return l.Command.Help(err.Error())
	}
	if args.delete {
		err = l.deleteLabels(args, c)
		if err != nil {
			l.Command.Help(err.Error())
		}
		return l.Command.TextMessage("Labels deleted")
	}
	err = l.storeLabels(args, c)
	if err != nil {
		l.Command.Help(err.Error())
	}

	return l.Command.TextMessage("Labels added")
}

func (l *Labeler) storeLabels(args args, c slack.SlashCommand) error {
	for _, label := range args.labels {
		err := l.storage.Store(c.ChannelID, label, commands.Label)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error save label to mapper")
		}
	}
	return nil
}

func (l *Labeler) deleteLabels(args args, c slack.SlashCommand) error {
	for _, label := range args.labels {
		err := l.storage.Delete(c.ChannelID, label, commands.Label)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error delete label from mapper")
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
	if pieces[0] == "-delete" {
		res.labels = pieces[1:]
		res.delete = true
	}
	return res, nil
}
