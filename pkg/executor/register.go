package executor

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type GitTagRegister struct {
	f                  makeGitTagMessage
	GitCommandExecutor GitCommandExecutor
}

type TagType string

const (
	Light     TagType = "light"
	Annotated TagType = "annotated"
)

type makeGitTagMessage func(*domain.CommitId, *domain.ServiceTagWithSemVer) string

func NewGitTagRegister(executor GitCommandExecutor, opt ...makeGitTagMessage) *GitTagRegister {
	f := func(commitId *domain.CommitId, tag *domain.ServiceTagWithSemVer) string {
		return fmt.Sprintf("Add %s tags to %s", tag.String(), commitId.String())
	}
	if len(opt) > 0 {
		f = opt[0]
	}
	return &GitTagRegister{
		f:                  f,
		GitCommandExecutor: executor,
	}
}

func (g *GitTagRegister) Execute(cmd usecase.RegisterServiceTagsCommand) error {
	for _, tag := range *cmd.Tags {
		if g.f == nil {
			_, err := gitTagAddLight(g.GitCommandExecutor, cmd.CommitId.String(), tag.String())
			if err != nil {
				return err
			}
			continue
		}
		message := g.f(cmd.CommitId, tag)
		_, err := gitTagAdd(g.GitCommandExecutor, cmd.CommitId.String(), tag.String(), message)
		if err != nil {
			return err
		}
	}
	return nil
}
