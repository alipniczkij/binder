package gitlab

type IssueInfo struct {
	ID       int
	Title    string
	Author   string
	Assignee string
	Labels   []string
}
