package gitlab

import (
	"github.com/kozaktomas/dashboard/pkg/integrations"
	"github.com/kozaktomas/dashboard/pkg/utils"
	"github.com/xanzy/go-gitlab"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Service struct {
	client     *gitlab.Client
	userId     int
	projects   []string
	itemsCache []integrations.Item
}

func New(token string, userId int, projects []string) integrations.Integration {
	client, _ := gitlab.NewClient(token)

	return &Service{
		client:   client,
		userId:   userId,
		projects: projects,
	}
}

func (gl *Service) GetName() string {
	return "gitlab"
}

func (gl *Service) GetItems() []integrations.Item {
	if len(gl.itemsCache) > 0 {
		return gl.itemsCache
	}

	var items []integrations.Item

	for _, projectId := range gl.projects {
		project, _, _ := gl.client.Projects.GetProject(projectId, &gitlab.GetProjectOptions{})
		authorId := gl.userId
		state := "opened"

		// author
		mrsAuthor, _, _ := gl.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			AuthorID: &authorId,
			State:    &state,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		// reviewer
		mrsReviewer, _, _ := gl.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			ReviewerID: &authorId,
			State:      &state,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		mrs := append(mrsAuthor, mrsReviewer...)
		sort.Slice(mrs[:], func(i, j int) bool {
			return mrs[i].IID > mrs[j].IID
		})

		for _, mr := range mrs {
			items = append(items, integrations.Item{
				Id:   strconv.Itoa(mr.IID) + "|" + projectId,
				Text: gl.formatGitlabMergeRequestTitle(project, mr),
				Url:  mr.WebURL,
				Copy: mr.SourceBranch,
			})
		}
	}

	gl.itemsCache = items
	return items
}

func (gl *Service) GetDetail(i integrations.Item) integrations.ItemDetail {
	splitId := strings.Split(i.Id, "|")
	iid, _ := strconv.Atoi(splitId[0])
	projectId, _ := strconv.Atoi(splitId[1])

	var wg sync.WaitGroup
	wg.Add(2)

	var mr *gitlab.MergeRequest
	go func() {
		defer wg.Done()
		mr, _, _ = gl.client.MergeRequests.GetMergeRequest(projectId, iid, &gitlab.GetMergeRequestsOptions{})
	}()

	var mrApprovals *gitlab.MergeRequestApprovals
	go func() {
		defer wg.Done()
		mrApprovals, _, _ = gl.client.MergeRequests.GetMergeRequestApprovals(projectId, iid)
	}()

	wg.Wait()

	assignees := make([]string, len(mr.Assignees))
	for i, ass := range mr.Assignees {
		assignees[i] = ass.Name
	}

	reviewers := make([]string, len(mr.Reviewers))
	for i, rew := range mr.Reviewers {
		reviewers[i] = rew.Name
	}

	conflict := "No"
	if mr.HasConflicts {
		conflict = "Yes"
	}

	approvedBy := make([]string, len(mrApprovals.ApprovedBy))
	for i, app := range mrApprovals.ApprovedBy {
		approvedBy[i] = app.User.Name
	}

	return integrations.ItemDetail{
		Title: "Who knows?",
		Parts: []integrations.Renderable{
			utils.Paragraph{Text: "[Title](fg:cyan): " + mr.Title},
			utils.Break{},
			utils.Break{},
			utils.Paragraph{Text: "[Description](fg:cyan):\n " + mr.Description},
			utils.Break{},
			utils.Break{},
			utils.Paragraph{Text: "[Assignees](fg:cyan): " + strings.Join(assignees, ", ")},
			utils.Break{},
			utils.Paragraph{Text: "[Approvals](fg:cyan): " + strings.Join(approvedBy, ", ")},
			utils.Break{},
			utils.Paragraph{Text: "[Reviewers](fg:cyan): " + strings.Join(reviewers, ", ")},
			utils.Break{},
			utils.Paragraph{Text: "[Status](fg:cyan): " + mr.State},
			utils.Break{},
			utils.Paragraph{Text: "[Branch](fg:cyan): " + mr.SourceBranch},
			utils.Break{},
			utils.Paragraph{Text: "[Conflict](fg:cyan): " + conflict},
			utils.Break{},
			utils.Paragraph{Text: "[Pipeline](fg:cyan): " + mr.Pipeline.Status},
			utils.Break{},
			utils.Paragraph{Text: "[Files changed](fg:cyan): " + mr.ChangesCount},
			utils.Break{},
			utils.Paragraph{Text: "[Created](fg:cyan): " + mr.CreatedAt.String()},
			utils.Break{},
		},
	}
}

func (gl *Service) formatGitlabMergeRequestTitle(project *gitlab.Project, mr *gitlab.MergeRequest) string {
	projectName := utils.RemoveWhiteSpaces(project.Namespace.Path + "/" + project.Name)
	projectName = strings.ToLower(projectName)

	badge := ""
	if mr.Author.ID == gl.userId {
		badge += "A"
	}

	if badge != "A" {
		for _, assignee := range mr.Assignees {
			if assignee.ID == gl.userId {
				badge += "A"
			}
		}
	}

	for _, reviewer := range mr.Reviewers {
		if reviewer.ID == gl.userId {
			badge += "[R](fg:red)"
		}
	}

	ret := " " + badge + " "
	ret += "[" + projectName + "!" + strconv.Itoa(mr.IID) + "](fg:cyan)"
	ret += " - " + mr.Title

	return ret
}
