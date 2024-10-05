package usecase

import (
	"msgtm/pkg/domain"
)

func SyncAllServiceTagState(state *domain.WritedState, list ListTags, finder CommitFinder) (*domain.WritedState, error) {
	tags, err := list.Execute(ListTagsQuery{
		Filter: func(_ *domain.ServiceName) bool {
			return true
		}})
	if err != nil {
		return nil, err
	}
	serviceTags := domain.FilterServiceTags(tags)
	sorts := domain.SortsServiceTags(serviceTags)
	for serviceName, tags := range sorts {
		latest := tags[len(tags)-1]
		gitTagLatest := latest.ToGitTag()
		query := FindCommitQuery{
			Tag: &gitTagLatest,
		}
		commitId, err := finder.Execute(query)
		if err != nil {
			return nil, err
		}
		info := domain.ServiceTagInfo{
			Tag:      latest,
			CommitId: commitId,
		}
		var prev *domain.ServiceTagInfo = nil
		if len(tags) > 1 {
			tag := tags[len(tags)-2]
			gitTag := tag.ToGitTag()
			commitId, err := finder.Execute(FindCommitQuery{
				Tag: &gitTag,
			})
			if err != nil {
				return nil, err
			}
			prev = &domain.ServiceTagInfo{
				Tag:      tag,
				CommitId: commitId,
			}
		}
		state.Update(serviceName, &info, prev)
	}
	return state, nil
}
