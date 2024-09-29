package subcmd

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type PushCommandParameter struct {
	CommitId string
	Remote   string
}

func PushCommand(getter usecase.CommitTagGetter, pusher usecase.CommitPusher) SubCommand[PushCommandParameter] {
	return func(param PushCommandParameter) error {
		commitId := domain.HEAD
		if param.CommitId != "" {
			commitId = domain.CommitId(param.CommitId)
		}
		remote := domain.Origin
		if param.Remote != "" {
			remote = domain.RemoteAddr(param.Remote)
		}

		err := usecase.PushAll(
			getter,
			pusher,
			&remote,
			&commitId,
		)
		if err != nil {
			return fmt.Errorf("failed to push service tags: %w", err)
		}
		return nil
	}
}
