package usecase

import "msgtm/pkg/domain"

func ResetServiceTags(destroyer DestroyServiceTags, commitGetter CommitTagGetter, commitId *domain.CommitId) error {
	tags, err := commitGetter.Execute(GetCommitTagQuery{CommitId: commitId})
	if err != nil {
		return err
	}
	targets := domain.FilterServiceTags(tags)
	err = destroyer.Execute(DestroyServiceTagsCommand{Tags: targets})
	if err != nil {
		return err
	}
	return nil
}
