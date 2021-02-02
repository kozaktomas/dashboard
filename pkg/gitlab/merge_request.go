package gitlab

import (
	"fmt"
	gl "github.com/xanzy/go-gitlab"
	"sort"
)

func IsMrOpen(mr *gl.MergeRequest) bool {
	if mr.State == "opened" {
		return true
	}

	return false
}

func (c *Service) GetMyMergeRequests() ([]*gl.MergeRequest, error) {
	var result []*gl.MergeRequest

	for _, projectId := range c.projects {
		authorId := c.userId
		mrs, _, err := c.client.MergeRequests.ListProjectMergeRequests(projectId, &gl.ListProjectMergeRequestsOptions{
			AuthorID: &authorId,
			ListOptions: gl.ListOptions{
				PerPage: 1000,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("could not retrieve merge request data from Gitlab: %w", err)
		}

		result = append(mrs, result...)
	}

	sort.Slice(result, func(i, j int) bool {
		// top prio for opened
		if result[i].State == "opened" && result[j].State != "opened" {
			return true
		}

		return false
	})

	return result, nil
}

