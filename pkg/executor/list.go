package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type GitTagList struct {
	GitCommandExecutor gitCommandExecutor
}

func (f *GitTagList) Execute(cmd usecase.ListTagsQuery) (*[]domain.GitTag, error) {
	tags, err := tagList(
		f.GitCommandExecutor,
	)
	if err != nil {
		return nil, err
	}
	filteredTags := []domain.GitTag{}
	for _, tag := range *tags {
		for _, service := range *cmd.Filter {
			if service.IsServiceTag(&tag) {
				filteredTags = append(filteredTags, tag)
			}
		}
	}
	return &filteredTags, nil
}
