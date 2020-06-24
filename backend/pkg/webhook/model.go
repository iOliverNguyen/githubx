package webhook

import "github.com/ng-vu/githubx/backend/pkg/model"

type WebhookRequest struct {
	Action      string              `json:"action"`
	Changes     *Changes            `json:"changes"`
	Issue       *Issue              `json:"issue"`
	Comment     *Comment            `json:"comment"`
	Repository  *model.RepositoryX  `json:"repository"`
	ProjectCard *model.ProjectCardX `json:"project_card"`
	Sender      *model.UserX        `json:"sender"`
}

type Issue struct {
	*model.IssueX

	User      *model.UserX      `json:"user"`
	Labels    model.LabelsX     `json:"labels"`
	Assignees model.UsersX      `json:"assignees"`
	Milestone *model.MilestoneX `json:"milestone"`
}

type Comment struct {
	*model.CommentX

	User *model.UserX `json:"user"`
}

type Changes struct {
	ColumnID ColumnChange `json:"column_id"`
}

type ColumnChange struct {
	From int `json:"from"`
}
