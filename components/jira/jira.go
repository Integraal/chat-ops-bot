package jira

import (
	client "github.com/andygrunwald/go-jira"
	"fmt"
	"errors"
	"github.com/integraal/chat-ops-bot/components/event"
	"github.com/integraal/chat-ops-bot/components/user"
	"encoding/json"
	"io/ioutil"
	"strings"
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
	epicKey     string
	epicField   string
}

type Config struct {
	Url         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Project     string `json:"project"`     // PROJECT
	EpicKey     string `json:"epicKey"`     // PROJECT-99
	EpicField   string `json:"epicField"`   // customfield_10008
	IssuePrefix string `json:"issuePrefix"` // [TimeTracking]
	IssueLabel  string `json:"issueLabel"`  // TimeTracking
	IssueType   string `json:"issueType"`   // Task
}

func Initialize(config Config) {
	cl, _ := client.NewClient(nil, config.Url)
	cl.Authentication.SetBasicAuth(config.Username, config.Password)
	//cl.Authentication.AcquireSessionCookie(config.Username, config.Password)
	jira = Jira{
		client:      cl,
		issuePrefix: config.IssuePrefix,
		issueLabel:  config.IssueLabel,
		project:     config.Project,
		username:    config.Username,
		issueType:   config.IssueType,
		epicField:   config.EpicField,
		epicKey:     config.EpicKey,
	}
}

func Get() *Jira {
	return &jira
}

func (j *Jira) EnsureIssue(event *event.Event) (*client.Issue, *client.Response, error) {
	var err error
	var response *client.Response
	issue := j.GetIssue(event.Uid)
	if issue == nil {
		issue, response, err = j.createIssue(event)
		if err != nil {
			return nil, response, err
		}
	}
	return issue, nil, nil
}

func (j *Jira) getIssueLabels(eventUid string) []string {
	return []string{
		fmt.Sprintf(j.issueLabel),
		fmt.Sprintf(j.issueLabel + ":Event:" + eventUid),
	}
}

func (j *Jira) getJQL(eventUid string) string {
	return fmt.Sprintf(JQLPattern, j.project, j.issueLabel+":Event:"+eventUid)
}

func (j *Jira) findIssue(eventUid string) (*client.Issue, error) {
	jql := j.getJQL(eventUid)
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

func (j *Jira) GetIssue(eventUid string) *client.Issue {
	issue, err := j.findIssue(eventUid)
	if err != nil {
		return nil
	} else {
		return issue
	}
}

func (j *Jira) getIssueSummary(event *event.Event) string {
	return j.issuePrefix + event.Summary
}

func (j *Jira) getIssueDescription(event *event.Event) string {
	return strings.Replace(event.Description, "\\n", "\n", -1)
}

func (j *Jira) createIssue(event *event.Event) (*client.Issue, *client.Response, error) {
	i := &client.Issue{
		Fields: &client.IssueFields{
			Project: client.Project{
				Key: j.project,
			},
			Type: client.IssueType{
				Name: j.issueType,
			},
			Summary:     j.getIssueSummary(event),
			Description: j.getIssueDescription(event),
			Labels:      j.getIssueLabels(event.Uid),
		},
	}

	issue, response, err := j.createIssueWithEpic(i)
	if err != nil {
		return nil, response, err
	}
	return issue, nil, nil
}

func (j *Jira) createIssueWithEpic(issue *client.Issue) (*client.Issue, *client.Response, error) {
	apiEndpoint := "rest/api/2/issue/"

	// Hack: add dynamic epic field via json marshalling
	m := make(map[string]map[string]interface{})
	b, _ := json.Marshal(issue)
	json.Unmarshal(b, &m)
	m["fields"][j.epicField] = j.epicKey
	// End of hack

	req, err := j.client.NewRequest("POST", apiEndpoint, m)
	if err != nil {
		return nil, nil, err
	}
	resp, err := j.client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	responseIssue := new(client.Issue)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	err = json.Unmarshal(data, responseIssue)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return responseIssue, resp, nil
}

type WorklogIssue struct {
	Key string `json:"key"`
}

type WorklogUser struct {
	Name string `json:"name"`
}

type Worklog struct {
	ID               int64 `json:"id,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Self             string `json:"self,omitempty"`
	Issue            *WorklogIssue `json:"issue,omitempty"`
	Author           *WorklogUser `json:"author,omitempty"`
	TimeSpentSeconds int64 `json:"timeSpentSeconds,omitempty"`
	BilledSeconds    int64 `json:"billedSeconds,omitempty"`
	DateStarted      string `json:"dateStarted,omitempty"`
}

func (j *Jira) AddUserTime(issue *client.Issue, evt *event.Event, user *user.User) (*client.Response, error) {
	worklog := Worklog{
		Comment:          "Присутствие на событии " + evt.Summary,
		TimeSpentSeconds: int64(evt.GetDuration().Seconds()),
		Author:           &WorklogUser{user.JiraUsername},
		Issue:            &WorklogIssue{issue.Key},
		DateStarted:      evt.GetStartTime().Format("2006-01-02T15:04:05.000"),
	}
	req, err := j.client.NewRequest("POST", "rest/tempo-timesheets/3/worklogs/", &worklog)

	if err != nil {
		return nil, err
	}

	resp, err := j.client.Do(req, nil)

	if err != nil {
		return resp, err
	}
	return resp, nil
}
