package msgtm

import (
	"fmt"
	"os/exec"
)

type GitTagRegister struct {
	f MakeGitTagMessage
}

type MakeGitTagMessage func(*CommitId, *ServiceTagWithSemVer) string

func DefaultGitTagRegister() *GitTagRegister {
	defaultFunc := func(commitId *CommitId, tag *ServiceTagWithSemVer) string {
		return fmt.Sprintf("Add %s tags to %s", tag.String(), commitId.String())
	}
	return &GitTagRegister{f: defaultFunc}
}

func (g *GitTagRegister) Register(commitId *CommitId, tags *[]*ServiceTagWithSemVer) error {
	for _, tag := range *tags {
		message := g.f(commitId, tag)
		_, err := gitTagAdd(commitId.String(), tag.String(), message)
		if err != nil {
			return err
		}
	}
	return nil
}

func gitTagAdd(commitId string, tag string, message string) (string, error) {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message, commitId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
