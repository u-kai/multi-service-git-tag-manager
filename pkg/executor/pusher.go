package executor

import "msgtm/pkg/usecase"

type GitTagPusher struct {
}

func (g *GitTagPusher) Execute(cmd usecase.CommitPushCommand) error {
	tagStrs := []string{}
	for _, tag := range *cmd.Tags {
		tagStrs = append(tagStrs, tag.String())
	}
	_, err := gitPushTags(cmd.RemoteAddr.String(), tagStrs...)
	if err != nil {
		return err
	}
	return nil
}
