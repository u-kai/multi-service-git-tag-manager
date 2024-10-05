package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type ServiceTagsListParameter struct {
	Filter []string
	IsAll  bool
}

func ServiceTagsListCommand(list usecase.ListTags, finder usecase.CommitFinder) SubCommand[ServiceTagsListParameter] {
	return func(param ServiceTagsListParameter) error {
		f := func(s *domain.ServiceName) bool {
			for _, filter := range param.Filter {
				if filter == s.String() {
					return true
				}
			}
			return false
		}
		if param.IsAll {
			f = func(s *domain.ServiceName) bool {
				return true
			}
		}
		infos, err := usecase.ServiceTagsList(f, list, finder)
		if err != nil {
			return fmt.Errorf("failed to list service tags: %w", err)
		}
		for _, info := range infos {
			fmt.Printf("%s:%s\n", info.Tag, info.CommitId)
		}
		return nil
	}
}
