package usecase

import "msgtm/pkg/domain"

func VersionUpAllServiceTags(
	list ListTags,
	registerService RegisterServiceTags,
	versionUpService domain.VersionUpServiceTag,
	commitId *domain.CommitId,
	excludeServiceNames ...*domain.ServiceName,
) error {
	f := func(s *domain.ServiceName) bool {
		for _, name := range excludeServiceNames {
			if *name == *s {
				return false
			}
		}
		return true
	}
	tags, err := list.Execute(ListTagsQuery{
		Filter: f,
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
