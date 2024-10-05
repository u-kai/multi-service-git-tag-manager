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
		serviceNames := []domain.ServiceName{}
		if param.FromConfigFile != "" {
			// read from config file
			file, err := os.Open(param.FromConfigFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			state, err := domain.FromReader(file, domain.YAML)
			for _, service := range state.ServiceTagStates {
				serviceNames = append(serviceNames, *service.ServiceName)
			}
		} else {
			for _, service := range param.Services {
				serviceNames = append(serviceNames, domain.ServiceName(service))
			}
		}

		commitId := domain.HEAD
		if param.CommitId != "" {
			commitId = domain.CommitId(param.CommitId)
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
