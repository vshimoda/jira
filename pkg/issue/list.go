package issue

import (
	"fmt"
	"strings"

	"github.com/3-shake/jira/pkg/auth"
	"github.com/3-shake/jira/pkg/project"
	"github.com/andygrunwald/go-jira"
)

// type Issue struct {
// ID     string  `yaml:"id"`
// Key    string  `yaml"key"`
// Fields *Fields `yaml:"fields"`
//
// fields *jira.IssueFields
// }
//
// func (i *Issue) Label() string {
// format := "%v: %v"
// return fmt.Sprintf(format, i.Key, i.Fields.Summary)
// }

// type Fields struct {
// Type        *jira.IssueType `yaml:"type,omitempty"`
// Labels      []string        `yaml:"labels,omitempty"`
// Summary     string          `yaml:"summary,omitempty"`
// Status      *jira.Status    `yaml:"status,omitempty"`
// Description string          `yaml:"description,omitempty"`
// Assignee    *jira.User      `yaml:"assignee,omitempty"`
// Reporter    *jira.User      `yaml:"reporter,omitempty"`
// }

type Search struct {
	Verbose         bool     `json:"-"`
	AffectedVersion string   `json:"AffectedVersion"`
	Assignee        string   `json:"assignee"`
	Comment         string   `json:"comment"`
	Component       string   `json:"component"`
	Created         string   `json:"Created"`
	Creator         string   `json:"Creator"`
	Description     string   `json:"description"`
	Due             string   `json:"due"`
	EpicLink        string   `json:"epic link"`
	Filter          string   `json:"filter"`
	IssueKey        string   `json:"issueKey"`
	Summary         string   `json:"summary"`
	Labels          []string `json:"labels"`
	Resolution      string   `json:"resolution"`
	Parent          string   `json:"parent"`
	Project         string   `json:"project"`
	Status          string   `json:"status"`
	Type            string   `json:"type"`
	Updated         string   `json:"updated"`
	Reporter        string   `json:"reporter"`
}

func List(op *Search) ([]jira.Issue, error) {
	projectName, err := project.Current()
	if err != nil {
		return nil, err
	}

	cli, err := auth.Client()
	if err != nil {
		return nil, err
	}
	jqlList := make([]string, 0)
	jqlList = append(jqlList, fmt.Sprintf("project = %v", projectName))
	if op.Summary != "" {
		jqlList = append(jqlList, fmt.Sprintf("summary ~ \"%v\"", op.Summary))
	}

	if op.Status != "" {
		jqlList = append(jqlList, fmt.Sprintf("status = %v", op.Status))
	}

	if len(op.Labels) > 0 {
		jqlList = append(jqlList, fmt.Sprintf("labels in (%v)", strings.Join(op.Labels, ",")))
	}

	if op.Assignee == "own" {
		a, err := auth.Read()
		if err != nil {
			return nil, err
		}
		jqlList = append(jqlList, fmt.Sprintf("assignee = \"%v\"", a.Username))
	} else if op.Assignee != "" {
		jqlList = append(jqlList, fmt.Sprintf("assignee = \"%v\"", op.Assignee))
	}

	if op.Reporter == "own" {
		a, err := auth.Read()
		if err != nil {
			return nil, err
		}
		jqlList = append(jqlList, fmt.Sprintf("reporter = \"%v\"", a.Username))
	} else if op.Reporter != "" {
		jqlList = append(jqlList, fmt.Sprintf("reporter = \"%v\"", op.Reporter))
	}

	jql := strings.Join(jqlList, " AND ")
	issues, _, err := cli.Issue.Search(jql, nil)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
