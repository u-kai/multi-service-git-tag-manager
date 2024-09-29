package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"strings"
)

type CommitFinder struct {
	GitCommandExecutor GitCommandExecutor
}

func (c *CommitFinder) Execute(query usecase.FindCommitQuery) (*domain.CommitId, error) {
	commitId, err := gitRevList(c.GitCommandExecutor, query.Tag.String())
	if err != nil {
		return nil, err
	}
	result := domain.CommitId(strings.Split(commitId, "\n")[0])
	return &result, nil
}
