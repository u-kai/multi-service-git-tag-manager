package usecase

import (
	"fmt"
	"msgtm/pkg/domain"
)

func SyncAllServiceTagState(state *domain.WritedState, list ListTags, finder CommitFinder) (*domain.WritedState, error) {
	tags, err := list.Execute(ListTagsQuery{})
	if err != nil {
		return nil, err
	}
	serviceTags := domain.FilterServiceTags(tags)
	sorts := domain.SortsServiceTags(serviceTags)
	for serviceName, tags := range sorts {
		latest := tags[0]
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
		fmt.Printf("serviceName: %v\n", serviceName)
		fmt.Printf("info: %v\n", info)
		state.Update(serviceName, &info)
		fmt.Printf("state: %v\n", state)
	}
	return state, nil
}
