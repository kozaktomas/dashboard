package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

type Service struct {
	client   *gitlab.Client
	userId   int
	projects []string
}

func New(key string, userId int, projectIds []string) (*Service, error) {
	client, err := gitlab.NewClient(key)
	if err != nil {
		return nil, err
	}
	return &Service{
		client:   client,
		userId:   userId,
		projects: projectIds,
	}, nil
}
