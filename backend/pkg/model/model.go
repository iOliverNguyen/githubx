package model

import "time"

var _ = UserX(User{})
var _ = IssueX(Issue{})
var _ = CommentX(Comment{})
var _ = LabelX(Label{})
var _ = MilestoneX(Milestone{})

type Issue struct {
	Body      string    `json:"body"`
	BodyHTML  string    `json:"bodyHTML"`
	BodyText  string    `json:"bodyText"`
	Closed    bool      `json:"closed"`
	ClosedAt  time.Time `json:"closedAt"`
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	Number    int       `json:"number"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
	URL       string    `json:"url"`
}

type IssueX struct {
	Body      string    `json:"body"`
	BodyHTML  string    `json:"body_html"`
	BodyText  string    `json:"body_text"`
	Closed    bool      `json:"closed"`
	ClosedAt  time.Time `json:"closed_at"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"node_id"`
	Number    int       `json:"number"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"html_url"`
}

type Author struct {
	Login string `json:"login"`
	URL   string `json:"url"`
}

type User struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatarURL"`
}

func (u *User) ToAuthor() Author {
	return Author{
		Login: u.Login,
		URL:   u.URL,
	}
}

type UserX struct {
	ID        string `json:"node_id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	URL       string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
}

func (u *UserX) ToUser() *User { return (*User)(u) }

type UsersX []*UserX

func (us UsersX) ToUsers() []*User {
	result := make([]*User, len(us))
	for i := range us {
		result[i] = (*User)(us[i])
	}
	return result
}

type Comment struct {
	Body      string    `json:"body"`
	BodyHTML  string    `json:"bodyHTML"`
	BodyText  string    `json:"bodyText"`
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentX struct {
	Body      string    `json:"body"`
	BodyHTML  string    `json:"body_html"`
	BodyText  string    `json:"body_text"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"node_id"`
	URL       string    `json:"html_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Label struct {
	Color       string `json:"color"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type LabelX struct {
	Color       string `json:"color"`
	Description string `json:"description"`
	ID          string `json:"node_id"`
	Name        string `json:"name"`
	URL         string `json:"html_url"`
}

type LabelsX []*LabelX

func (us LabelsX) ToLabels() []*Label {
	result := make([]*Label, len(us))
	for i := range us {
		result[i] = (*Label)(us[i])
	}
	return result
}

type RepositoryX struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type Milestone struct {
	ClosedAt     time.Time `json:"closed_at"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	ID           string    `json:"id"`
	Number       int       `json:"number"`
	OpenIssues   int       `json:"open_issues"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	UpdatedAt    time.Time `json:"updated_at"`
	URL          string    `json:"url"`
}

type MilestoneX struct {
	ClosedAt     time.Time `json:"closed_at"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	ID           string    `json:"node_id"`
	Number       int       `json:"number"`
	OpenIssues   int       `json:"open_issues"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	UpdatedAt    time.Time `json:"updated_at"`
	URL          string    `json:"html_url"`
}

type ProjectCardX struct {
	Archived   bool   `json:"archived"`
	ColumnID   int    `json:"column_id"`
	ContentURL string `json:"content_url"`
	ID         string `json:"node_id"`
	URL        string `json:"html_url"`
}
