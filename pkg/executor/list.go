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
		serviceTag, err := tag.ToServiceTag()
		if err != nil {
			continue
		}
		if cmd.Filter(&serviceTag.Service) {
			filteredTags = append(filteredTags, tag)
		}
	}
	return &filteredTags, nil
}
