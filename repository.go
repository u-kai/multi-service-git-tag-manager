package msgtm

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitTagRegister struct {
	f       MakeGitTagMessage
	handler EventHandler[RegisterEvent]
}

type EventHandler[T any] func(T) error

type TagType string

const (
	Light     TagType = "light"
	Annotated TagType = "annotated"
)

type RegisterEvent struct {
	Type     TagType
	CommitId *CommitId
	Tag      *ServiceTagWithSemVer
}

func logRegisterEvent(event RegisterEvent) error {
	fmt.Printf("Register %s tag %s to %s\n", event.Type, event.Tag.String(), event.CommitId.String())
	return nil
}

type MakeGitTagMessage func(*CommitId, *ServiceTagWithSemVer) string

func GitTagRegisterWithDefaultMessage() *GitTagRegister {
	defaultFunc := func(commitId *CommitId, tag *ServiceTagWithSemVer) string {
		return fmt.Sprintf("Add %s tags to %s", tag.String(), commitId.String())
	}
	return &GitTagRegister{
		f:       defaultFunc,
		handler: logRegisterEvent,
	}
}
func DefaultGitTagRegister() *GitTagRegister {
	return &GitTagRegister{
		handler: logRegisterEvent,
	}
}

func (g *GitTagRegister) Register(commitId *CommitId, tags *[]*ServiceTagWithSemVer) error {
	for _, tag := range *tags {
		if g.f == nil {
			_, err := gitTagAddLight(commitId.String(), tag.String())
			if err != nil {
				return err
			}
			g.handler(RegisterEvent{
				Type:     Light,
				CommitId: commitId,
				Tag:      tag,
			})
			continue
		}
		message := g.f(commitId, tag)
		_, err := gitTagAdd(commitId.String(), tag.String(), message)
		if err != nil {
			return err
		}
		g.handler(RegisterEvent{
			Type:     Annotated,
			CommitId: commitId,
			Tag:      tag,
		})

	}
	return nil
}

type AllTagList struct{}

func (a *AllTagList) List() (*[]GitTag, error) {
	return tagList()
}

type FilterTagList struct {
	IncludePrefix []string
}

func (f *FilterTagList) List() (*[]GitTag, error) {
	tags, err := tagList()
	if err != nil {
		return nil, err
	}
	filteredTags := []GitTag{}
	for _, tag := range *tags {
		for _, prefix := range f.IncludePrefix {
			if strings.HasPrefix(string(tag), prefix) {
				filteredTags = append(filteredTags, tag)
			}
		}
	}
	return &filteredTags, nil
}

func tagList() (*[]GitTag, error) {
	cmd := exec.Command("git", "tag")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(output), "\n")
	tagList := []GitTag{}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		tagList = append(tagList, GitTag(tag))
	}

	return &tagList, nil
}

func gitTagAddLight(commitId string, tag string) (string, error) {
	cmd := exec.Command("git", "tag", tag, commitId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
func gitTagAdd(commitId string, tag string, message string) (string, error) {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message, commitId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
