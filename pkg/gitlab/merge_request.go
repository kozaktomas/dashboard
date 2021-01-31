package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"sort"
	"time"
)

type MergeRequest struct {
	Id          int
	ProjectName string
	Title       string
	Url         string
	State       string
	Author      string
	AuthorImage string
	BranchName  string
}

func (mr *MergeRequest) IsOpen() bool {
	if mr.State == "opened" {
		return true
	}

	return false
}

func (mr *MergeRequest) IsClosed() bool {
	if mr.State == "closed" {
		return true
	}

	return false
}

func (mr *MergeRequest) IsMerged() bool {
	if mr.State == "merged" {
		return true
	}

	return false
}

func (c *Service) GetMyMergeRequests() ([]MergeRequest, error) {
	var mrs []MergeRequest

	for _, projectId := range c.projects {

		project, _, err := c.client.Projects.GetProject(projectId, &gitlab.GetProjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not retrieve project data from Gitlab: %w", err)
		}

		authorId := c.userId
		result, _, err := c.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			AuthorID: &authorId,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("could not retrieve merge request data from Gitlab: %w", err)
		}

		for _, mr := range result {
			mrs = append(mrs, MergeRequest{
				Id:          mr.IID,
				ProjectName: project.NameWithNamespace,
				Title:       mr.Title,
				Url:         mr.WebURL,
				State:       mr.State,
				Author:      mr.Author.Name,
				AuthorImage: mr.Author.AvatarURL,
				BranchName:  mr.SourceBranch,
			})
		}

		sort.Slice(mrs, func(i, j int) bool {
			// top prio for opened
			if mrs[i].IsOpen() && !mrs[j].IsOpen() {
				return true
			}

			return false
		})
	}

	return mrs, nil
}

// todo filter by userId
func (c *Service) GetReviews() ([]MergeRequest, error) {
	var mrs []MergeRequest

	for _, projectId := range c.projects {

		project, _, err := c.client.Projects.GetProject(projectId, &gitlab.GetProjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not retrieve project data from Gitlab: %w", err)
		}

		monthAgo := time.Now().AddDate(0, -2, 0)
		state := "opened"
		result, _, err := c.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			State:        &state,
			CreatedAfter: &monthAgo,
			ListOptions: gitlab.ListOptions{
				PerPage: 1000,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("could not retrieve merge request data from Gitlab: %w", err)
		}

		for _, mr := range result {
			mrs = append(mrs, MergeRequest{
				Id:          mr.IID,
				ProjectName: project.NameWithNamespace,
				Title:       mr.Title,
				Url:         mr.WebURL,
				State:       mr.State,
				Author:      mr.Author.Name,
				AuthorImage: mr.Author.AvatarURL,
				BranchName:  mr.SourceBranch,
			})
		}

	}

	return mrs, nil
}
