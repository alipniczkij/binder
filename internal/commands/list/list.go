package list

import (
	"fmt"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/storage"
	"github.com/slack-go/slack"
	"strings"
)

type List struct {
	commands.Command
	storage storage.Mapper
}

func New(s storage.Mapper) *List {
	return &List{
		Command: commands.Command{
			Usage: "{command} (ex. subscribe)",
		},
		storage: s,
	}
}

type args struct {
	command *string
}

func (l *List) Execute(c slack.SlashCommand) slack.Msg {
	args, err := l.parseArgs(c.Text)
	if err != nil {
		return l.Command.Help(err.Error())
	}
	msg, err := l.processList(*args.command, c.ChannelID)
	if err != nil {
		return l.Command.Help(err.Error())
	}
	return l.TextMessage(msg)
}

func (l *List) processList(cmd, channelID string) (string, error) {
	projects := make(map[string]struct{})

	keys, err := l.storage.GetKeys(cmd)
	if err != nil {
		return "", err
	}

	for _, k := range keys {
		channels, found := l.storage.Get(k, cmd)
		if found {
			for _, c := range channels {
				if c == channelID {
					projects[k] = struct{}{}
				}
			}
		}
	}

	projectsSlice := make([]string, len(projects))
	for k := range projects {
		projectsSlice = append(projectsSlice, k)
	}

	msg := fmt.Sprintf("There is projects finded: %s", strings.Join(projectsSlice, ", "))
	return msg, nil
}

func (l *List) parseArgs(input string) (args, error) {
	pieces := l.Command.SplitCommandLine(input)
	res := args{}
	if len(pieces) != 1 {
		return res, fmt.Errorf("invalid number of arguments")
	}
	res.command = &pieces[0]
	return res, nil
}
