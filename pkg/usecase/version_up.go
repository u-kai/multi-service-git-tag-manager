package usecase

import "msgtm/pkg/domain"

func VersionUpAllServiceTags(
	list ListTags,
	registerService RegisterServiceTags,
	versionUpService domain.VersionUpServiceTag,
	commitId *domain.CommitId,
	excludeServiceNames ...*domain.ServiceName,
) error {
	filter := make([]*domain.ServiceName, 0, len(excludeServiceNames))
	for _, name := range excludeServiceNames {
		filter = append(filter, name)
	}
	tags, err := list.Execute(ListTagsQuery{
		Filter: &filter,
	})
	if err != nil {
		return err
	}

	updates := versionUpService(tags)

	err = registerService.Execute(RegisterServiceTagsCommand{
		CommitId: commitId,
		Tags:     updates,
	})
	if err != nil {
		return err
	}

	return nil
}
