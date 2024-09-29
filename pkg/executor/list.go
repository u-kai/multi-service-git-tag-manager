package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type GitTagList struct {
	GitCommandExecutor GitCommandExecutor
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
		if cmd.Filter == nil {
			filteredTags = append(filteredTags, tag)
			continue
		}
		for _, filter := range *cmd.Filter {
			if filter.IsServiceTag(&tag) {
				filteredTags = append(filteredTags, tag)
			}
		}
	}
	return &filteredTags, nil
}
