package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type DestroyDecorator struct {
	Clients []usecase.DestroyServiceTags
}

func (d *DestroyDecorator) Execute(cmd usecase.DestroyServiceTagsCommand) error {
	for _, client := range d.Clients {
		err := client.Execute(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResetCommandParameter struct {
	Origin       bool
	ExcludeLocal bool
	CommitId     string
}

func ResetCommand(getter usecase.CommitTagGetter, local usecase.DestroyServiceTags, remote usecase.DestroyServiceTags) SubCommand[ResetCommandParameter] {
	return func(param ResetCommandParameter) error {
		commitId := domain.HEAD
		if param.CommitId != "" {
			commitId = domain.CommitId(param.CommitId)
		}

		destroyer := &DestroyDecorator{}
		if param.Origin {
			destroyer.Clients = append(destroyer.Clients, remote)
		}
		if !param.ExcludeLocal {
			destroyer.Clients = append(destroyer.Clients, local)
		}

		err := usecase.ResetServiceTags(
			destroyer,
			getter,
			&commitId,
		)
		if err != nil {
			return fmt.Errorf("failed to reset service tags: %w", err)
		}
		return nil
	}
}
