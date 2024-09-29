package executor

import (
	"fmt"
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type GitTagRegister struct {
	f makeGitTagMessage
}

type TagType string

const (
	Light     TagType = "light"
	Annotated TagType = "annotated"
)

type makeGitTagMessage func(*domain.CommitId, *domain.ServiceTagWithSemVer) string

func NewGitTagRegister(opt ...makeGitTagMessage) *GitTagRegister {
	f := func(commitId *domain.CommitId, tag *domain.ServiceTagWithSemVer) string {
		return fmt.Sprintf("Add %s tags to %s", tag.String(), commitId.String())
	}
	if len(opt) > 0 {
		f = opt[0]
	}
	return &GitTagRegister{
		f: f,
	}
}

func (g *GitTagRegister) Execute(cmd usecase.RegisterServiceTagsCommand) error {
	for _, tag := range *cmd.Tags {
		if g.f == nil {
			_, err := gitTagAddLight(cmd.CommitId.String(), tag.String())
			if err != nil {
				return err
			}
			continue
		}
		message := g.f(cmd.CommitId, tag)
		_, err := gitTagAdd(cmd.CommitId.String(), tag.String(), message)
		if err != nil {
			return err
		}
	}
	return nil
}
