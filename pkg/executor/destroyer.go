package executor

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type LocalServiceTagsDestroyer struct {
	Force              bool
	GitCommandExecutor GitCommandExecutor
}

func (s *LocalServiceTagsDestroyer) Execute(cmd usecase.DestroyServiceTagsCommand) error {
	for _, tag := range *cmd.Tags {
		_, err := gitTagDelete(s.GitCommandExecutor, tag.String(), s.Force)
		if err != nil {
			return err
		}
	}
	return nil
}

type RemoteServiceTagsDestroyer struct {
	Force              bool
	Remote             *domain.RemoteAddr
	GitCommandExecutor GitCommandExecutor
}

func (r *RemoteServiceTagsDestroyer) Execute(cmd usecase.DestroyServiceTagsCommand) error {
	tagStrs := []string{}
	for _, tag := range *cmd.Tags {
		tagStrs = append(tagStrs, tag.String())
	}
	_, err := gitTagRemoteDelete(r.GitCommandExecutor, r.Remote.String(), tagStrs)
	if err != nil {
		return err
	}
	return nil
}
