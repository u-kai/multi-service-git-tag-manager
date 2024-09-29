package usecase

import "msgtm/pkg/domain"

func CreateServiceTags(
	registerService RegisterServiceTags,
	commitId *domain.CommitId,
	serviceNames []domain.ServiceName,
	version domain.SemVer,
) error {
	serviceTags := make([]*domain.ServiceTagWithSemVer, 0, len(serviceNames))
	for _, serviceName := range serviceNames {
		serviceTags = append(serviceTags, domain.NewServiceTagWithSemVer(serviceName, version))
	}
	return registerService.Execute(RegisterServiceTagsCommand{
		CommitId: commitId,
		Tags:     &serviceTags,
	})
}
