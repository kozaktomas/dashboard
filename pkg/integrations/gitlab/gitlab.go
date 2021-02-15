package gitlab

import (
	"github.com/kozaktomas/dashboard/pkg/integrations"
	"github.com/kozaktomas/dashboard/pkg/utils"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"strings"
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
			State: &state,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		// reviewer
		mrsReviewer, _, _ := gl.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			ReviewerID: &authorId,
			State: &state,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		mrs := append(mrsAuthor, mrsReviewer...)

		for _, mr := range mrs {
			items = append(items, integrations.Item{
				Id:   strconv.Itoa(mr.IID) + "|" + projectId,
				Text: formatGitlabMergeRequestTitle(project, mr),
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

	mr, _, _ := gl.client.MergeRequests.GetMergeRequest(projectId, iid, &gitlab.GetMergeRequestsOptions{})

	assignees := make([]string, len(mr.Assignees))
	for i, ass := range mr.Assignees {
		assignees[i] = ass.Name
	}

	conflict := "No"
	if mr.HasConflicts {
		conflict = "Yes"
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
			utils.Paragraph{Text: "[Status](fg:cyan): " + mr.State},
			utils.Break{},
			utils.Paragraph{Text: "[Approvals](fg:cyan): " + strconv.Itoa(mr.ApprovalsBeforeMerge)},
			utils.Break{},
			utils.Paragraph{Text: "[Branch](fg:cyan): " + mr.SourceBranch},
			utils.Break{},
			utils.Paragraph{Text: "[Conflict](fg:cyan): " + conflict},
			utils.Break{},
			utils.Paragraph{Text: "[Pipeline](fg:cyan): " + mr.Pipeline.Status},
			utils.Break{},
			utils.Paragraph{Text: "[Files changed](fg:cyan): " + mr.ChangesCount},
			utils.Break{},
		},
	}
}

func formatGitlabMergeRequestTitle(project *gitlab.Project, mr *gitlab.MergeRequest) string {
	projectName := utils.RemoveWhiteSpaces(project.Namespace.Path + "/" + project.Name)
	projectName = strings.ToLower(projectName)

	ret := " "
	ret += "[" + projectName + "!" + strconv.Itoa(mr.IID) + "](fg:cyan)"
	ret += " - " + mr.Title

	return ret
}
