package models

type GitlabEvent struct {
	ObjectKind       string     `json:"object_kind"`
	EventType        string     `json:"event_type"`
	User             User       `json:"user"`
	Project          Project    `json:"project"`
	ObjectAttributes Attributes `json:"object_attributes"`
	Repository       Repository `json:"repository"`
	Assignees        []Assignee `json:"assignees"`
	Assignee         Assignee   `json:"assignee"`
	Labels           []Label    `json:"labels"`
}

type User struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type Project struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WebURL      string `json:"web_url"`
	AvatarURL   string `json:"avatar_url"`
}

type Attributes struct {
	ID          uint64  `json:"id"`
	Title       string  `json:"title"`
	AssigneeIDs []int   `json:"assignee_ids"`
	AssigneeID  int     `json:"assignee_id"`
	AuthorID    int     `json:"author_id"`
	ProjectID   int     `json:"project_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	State       string  `json:"state"`
	Action      string  `json:"action"`
	Labels      []Label `json:"labels"`
}

type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

type Assignee struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type Label struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	ProjectID   int    `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int    `json:"group_id"`
}
