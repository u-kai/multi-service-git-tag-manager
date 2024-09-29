package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type ServiceTagsListParameter struct {
	Filter []string
}

func ServiceTagsListCommand(list usecase.ListTags, finder usecase.CommitFinder) SubCommand[ServiceTagsListParameter] {
	return func(param ServiceTagsListParameter) error {
		serviceNames := make([]*domain.ServiceName, 0, len(param.Filter))
		for _, name := range param.Filter {
			service := domain.ServiceName(name)
			serviceNames = append(serviceNames, &service)
		}
		infos, err := usecase.ServiceTagsList(&serviceNames, list, finder)
		if err != nil {
			return fmt.Errorf("failed to list service tags: %w", err)
		}
		for _, info := range infos {
			fmt.Printf("%s:%s\n", info.Tag, info.CommitId)
		}
		return nil
	}
}
