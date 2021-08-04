package commands

import (
	"fmt"
	"github.com/google/shlex"
	"github.com/slack-go/slack"
	"strings"
	"unicode"
)

type Handler interface {
	Execute(slack.SlashCommand) slack.Msg
}

const (
	Subscribe   = "/subscribe"
	Unsubscribe = "/unsubscribe"
	Label       = "/label"
	Unlabel     = "/unlabel"
	List        = "/list"
)

type Command struct {
	Usage string
}

func (me Command) TrimCommandLine(input string) string {
	return strings.TrimFunc(
		input,
		func(r rune) bool {
			//so-called object-replacement char. Appears at new lines for some reason
			return r == '\uFFFC' || unicode.IsSpace(r)
		},
	)
}

func (me Command) SplitCommandLine(input string) []string {
	clean := me.TrimCommandLine(input)
	pieces, _ := shlex.Split(clean)
	return pieces
}

func (me Command) TextMessage(text string) slack.Msg {
	return slack.Msg{
		ResponseType: slack.ResponseTypeInChannel,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.SectionBlock{
					Type: "section",
					Text: slack.NewTextBlockObject(
						slack.MarkdownType,
						text,
						false,
						false,
					),
				},
			},
		},
	}
}

func (me Command) Help(error string) slack.Msg {
	return slack.Msg{
		ResponseType: slack.ResponseTypeInChannel,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.SectionBlock{
					Type: "section",
					Text: slack.NewTextBlockObject(
						slack.MarkdownType,
						fmt.Sprintf("Invalid command:\n%v", error),
						false,
						false,
					),
				},
				slack.SectionBlock{
					Type: "section",
					Text: slack.NewTextBlockObject(
						slack.MarkdownType,
						fmt.Sprintf("*USAGE:*\n```%v```", me.Usage),
						false,
						false,
					),
				},
			},
		},
	}
}
