package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"strings"
)

type CommitTagGetter struct {
}

func (c *CommitTagGetter) Execute(query usecase.GetCommitTagQuery) (*[]domain.GitTag, error) {
	tags := []domain.GitTag{}
	output, err := gitShowCommit(query.CommitId.String())
	if err != nil {
		return nil, err
	}
	commitLine := strings.Split(output, "\n")[0]
	// tag: service1-v1.1.1, tags: service2-v1.1.1
	tagsStr := strings.Split(commitLine, "(")[1]
	// remove ")"
	tagsStr = tagsStr[:len(tagsStr)-1]

	for _, tagStr := range strings.Split(tagsStr, ", ") {
		if !strings.HasPrefix(tagStr, "tag: ") {
			continue
		}
		tag := domain.GitTag(
			strings.Split(tagStr, "tag: ")[1],
		)
		tags = append(tags, tag)
	}
	return &tags, nil
}
