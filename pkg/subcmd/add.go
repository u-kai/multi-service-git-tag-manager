package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"os"
)

type TagAddCommandParameter struct {
	Version        string
	CommitId       string
	Services       []string
	FromConfigFile string
}

func TagAddCommand(register usecase.RegisterServiceTags) SubCommand[TagAddCommandParameter] {
	return func(param TagAddCommandParameter) error {
		semVer, err := domain.FromStr(param.Version)
		if err != nil {
			return fmt.Errorf("failed to parse version: %w", err)
		}
		if param.FromConfigFile != "" {
			// read from config file
			content, err := os.ReadFile(param.FromConfigFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			// TODO
			println(content)
		}
		commitId := domain.HEAD
		if param.CommitId != "" {
			commitId = domain.CommitId(param.CommitId)
		}
		serviceNames := []domain.ServiceName{}
		for _, service := range param.Services {
			serviceNames = append(serviceNames, domain.ServiceName(service))
		}
		err = usecase.CreateServiceTags(
			register,
			&commitId,
			serviceNames,
			semVer,
		)
		if err != nil {
			return fmt.Errorf("failed to create service tags: %w", err)
		}
		return nil
	}
}
