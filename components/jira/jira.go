package jira

import (
	client "github.com/andygrunwald/go-jira"
	"fmt"
	"errors"
	"github.com/integraal/chat-ops-bot/components/event"
	"strconv"
)

var jira Jira

const JQLPattern = "project = %s and labels = %s"

type Jira struct {
	client      *client.Client
	username    string
	project     string
	issuePrefix string
	issueLabel  string
	issueType   string
}

type JiraConfig struct {
	Url         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Project     string `json:"project"`     // GPS
	IssuePrefix string `json:"issuePrefix"` // [TimeTracking]
	IssueLabel  string `json:"issueLabel"`  // TimeTracking
	IssueType   string `json:"issueType"`   // Task
}

func Initialize(config JiraConfig) {
	cl, _ := client.NewClient(nil, "https://jira.atlassian.com/")
	cl.Authentication.SetBasicAuth(config.Username, config.Password)
	jira = Jira{
		client:      cl,
		issuePrefix: config.IssuePrefix,
		issueLabel:  config.IssueLabel,
		project:     config.Project,
		username:    config.Username,
		issueType:   config.IssueType,
	}
}

func (j *Jira) EnsureIssue(event *event.Event) (*client.Issue, error) {
	var err error
	issue := j.GetIssue(event.ID)
	if issue == nil {
		issue, err = j.createIssue(event)
		if err != nil {
			return nil, err
		}
	}
	return issue, nil
}

func (j *Jira) getIssueLabels(eventId int64) []string {
	return []string{
		fmt.Sprintf(j.issueLabel),
		fmt.Sprintf(j.issueLabel + ":Event:" + strconv.FormatInt(eventId, 10)),
	}
}

func (j *Jira) getJQL(eventId int64) string {
	return fmt.Sprintf(JQLPattern, j.project, j.issueLabel+":Event:"+strconv.FormatInt(eventId, 10))
}

func (j *Jira) findIssue(eventId int64) (*client.Issue, error) {
	jql := j.getJQL(eventId)
	options := client.SearchOptions{
		MaxResults: 1,
	}
	issues, _, err := j.client.Issue.Search(jql, &options)
	if err != nil {
		return nil, err
	} else if len(issues) > 0 {
		return &issues[0], nil
	} else {
		return nil, errors.New("Issue does not exist")
	}
}

func (j *Jira) GetIssue(eventId int64) *client.Issue {
	issue, err := j.findIssue(eventId)
	if err != nil {
		return nil
	} else {
		return issue
	}
}

func (j *Jira) getIssueSummary(event *event.Event) string {
	return j.issuePrefix + event.Name
}

func (j *Jira) getIssueDescription(event *event.Event) string {
	return event.Description
}

func (j *Jira) createIssue(event *event.Event) (*client.Issue, error) {
	issue, _, err := j.client.Issue.Create(&client.Issue{
		Fields: &client.IssueFields{
			Project: client.Project{
				Key: j.project,
			},
			Type: client.IssueType{
				Name: j.issueType,
			},
			Summary:     j.getIssueSummary(event),
			Description: j.getIssueDescription(event),
			Labels:      j.getIssueLabels(event.ID),
			Reporter: &client.User{
				Name: j.username,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return issue, nil
}