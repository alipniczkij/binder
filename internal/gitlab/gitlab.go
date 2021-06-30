package gitlab

import (
	"errors"
	"fmt"
	"github.com/alipniczkij/binder/pkg/config"
	"github.com/slack-go/slack"
	g "github.com/xanzy/go-gitlab"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	client *g.Client
}

func New(cfg config.Gitlab) (*Client, error) {
	client, err := g.NewClient(cfg.Token, g.WithBaseURL(cfg.BaseURL))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c Client) Issue(project string, id int) (IssueInfo, error) {
	i, _, err := c.client.Issues.GetIssue(project, id)
	if err != nil {
		return IssueInfo{}, err
	}
	assigneeName := ""
	if i.Assignee != nil {
		assigneeName = i.Assignee.Name
	}
	info := IssueInfo{
		ID:       i.IID,
		Title:    i.Title,
		Author:   i.Author.Name,
		Labels:   i.Labels,
		Assignee: assigneeName,
	}
	return info, nil
}

func (c Client) Unfurl(u *url.URL) (slack.Attachment, error) {
	path := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(path) < 2 {
		return slack.Attachment{}, errors.New("no unfurl found")
	}
	if path[len(path)-2] != "issues" {
		return slack.Attachment{}, errors.New("no unfurl found")
	}
	issueId, err := strconv.ParseInt(path[len(path)-1], 10, 64)
	if err != nil {
		return slack.Attachment{}, errors.New("invalid issue id")
	}
	projectID, err := GetProjectIDFromURL(u)
	if err != nil {
		return slack.Attachment{}, err
	}
	info, err := c.Issue(projectID, int(issueId))
	if err != nil {
		return slack.Attachment{}, fmt.Errorf("no issue info found. projectID: %s. issueID: %s. err: %s", projectID, path[len(path)-1], err.Error())
	}
	if info.Assignee == "" {
		info.Assignee = "--"
	}
	labelsLine := ""
	if len(info.Labels) > 0 {
		labelsLine = "\n`" + strings.Join(info.Labels, "` `") + "`"
	}
	return slack.Attachment{
		Color: "#E67E22",
		Title: fmt.Sprintf("[#%v] %v", info.ID, info.Title),
		Text: fmt.Sprintf(
			"Reporter: %v\nAssignee: %v%v",
			info.Author,
			info.Assignee,
			labelsLine,
		),
	}, nil
}

func GetProjectIDFromURL(u *url.URL) (string, error) {
	path := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(path) < 2 {
		return "", errors.New("project ID not found")
	}
	projectIdPieces := path[:len(path)-2]
	if projectIdPieces[len(projectIdPieces)-1] == "-" {
		projectIdPieces = projectIdPieces[:len(projectIdPieces)-1]
	}
	projectId := strings.Join(projectIdPieces, "/")
	if projectId == "" {
		return "", errors.New("project ID not found")
	}
	return projectId, nil
}
