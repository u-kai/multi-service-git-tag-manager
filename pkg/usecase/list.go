package usecase

import "msgtm/pkg/domain"

type ServiceTagInfo struct {
	Tag      *domain.ServiceTagWithSemVer
	CommitId *domain.CommitId
}

func ServiceTagsList(filter *[]*domain.ServiceName, list ListTags, finder CommitFinder) ([]*ServiceTagInfo, error) {
	tags, err := list.Execute(ListTagsQuery{Filter: filter})
	if err != nil {
		return nil, err
	}
	serviceTags := domain.FilterServiceTags(tags)

	infos := make([]*ServiceTagInfo, 0, len(*serviceTags))

	for _, tag := range *serviceTags {
		gitTag := tag.ToGitTag()
		commitId, err := finder.Execute(FindCommitQuery{Tag: &gitTag})
		if err != nil {
			return nil, err
		}
		infos = append(infos, &ServiceTagInfo{
			Tag:      tag,
			CommitId: commitId,
		})
	}
	return infos, nil
}
