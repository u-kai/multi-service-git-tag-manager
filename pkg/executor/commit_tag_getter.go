package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type CommitTagGetter struct {
	GitCommandExecutor GitCommandExecutor
}

func (c *CommitTagGetter) Execute(query usecase.GetCommitTagQuery) (*[]domain.GitTag, error) {
	result := []domain.GitTag{}
	tags, err := gitShowCommitTags(c.GitCommandExecutor, query.CommitId.String())
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		result = append(result, domain.GitTag(tag))
	}
	return &result, nil
}
