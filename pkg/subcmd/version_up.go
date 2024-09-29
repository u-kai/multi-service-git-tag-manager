package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type VersionUpCommandParameter struct {
	Minor    bool
	Major    bool
	IsAll    bool
	CommitId string
	Services []string
}

func VersionUpCommand(list usecase.ListTags, register usecase.RegisterServiceTags, getter usecase.CommitTagGetter) SubCommand[VersionUpCommandParameter] {
	return func(param VersionUpCommandParameter) error {
		commitId := domain.HEAD
		if param.CommitId != "" {
			commitId = domain.CommitId(param.CommitId)
		}

		excludeServices := []*domain.ServiceName{}
		if !param.IsAll && len(param.Services) > 0 {
			for _, service := range param.Services {
				s := domain.ServiceName(service)
				excludeServices = append(excludeServices, &s)
			}
		}

		f := domain.PatchUpAll
		if param.Minor {
			f = domain.MinorUpAll
		}
		if param.Major {
			f = domain.MajorUpAll
		}

		err := usecase.VersionUpAllServiceTags(
			list,
			register,
			f,
			&commitId,
			excludeServices...,
		)
		if err != nil {
			return fmt.Errorf("failed to version up: %w", err)
		}
		return nil
	}
}
