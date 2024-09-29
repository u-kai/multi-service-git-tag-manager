package executor

import "msgtm/pkg/usecase"

type GitTagPusher struct {
	GitCommandExecutor gitCommandExecutor
}

func (g *GitTagPusher) Execute(cmd usecase.CommitPushCommand) error {
	tagStrs := []string{}
	for _, tag := range *cmd.Tags {
		tagStrs = append(tagStrs, tag.String())
	}
	_, err := gitPushTags(g.GitCommandExecutor, cmd.RemoteAddr.String(), tagStrs...)
	if err != nil {
		return err
	}
	return nil
}
