package gitlab

import (
	"fmt"
	gl "github.com/xanzy/go-gitlab"
)

var projectsCache []*gl.Project

func (c *Service) GetProject(projectId int) (*gl.Project, error) {
	for _, projectCache := range projectsCache {
		if projectId == projectCache.ID {
			return projectCache, nil
		}
	}

	project, _, err := c.client.Projects.GetProject(projectId, &gl.GetProjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve project data from Gitlab: %w", err)
	}

	projectsCache = append(projectsCache, project)
	return project, nil
}
