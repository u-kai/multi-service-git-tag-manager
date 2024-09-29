package executor

import (
	"msgtm/pkg/domain"
	"os/exec"
	"strings"
)

func tagList() (*[]domain.GitTag, error) {
	cmd := exec.Command("git", "tag")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(output), "\n")
	tagList := []domain.GitTag{}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		tagList = append(tagList, domain.GitTag(tag))
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

func gitTagDelete(tag string, force bool) (string, error) {
	deleteOption := "-d"
	if force {
		deleteOption = "-d"
	}
	cmd := exec.Command("git", "tag", deleteOption, tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func gitTagRemoteDelete(remote string, tags []string) (string, error) {
	deleteOption := "--delete"
	cmdArgs := []string{"push", remote, deleteOption}
	cmdArgs = append(cmdArgs, tags...)
	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func gitShowCommit(commitId string) (string, error) {
	cmd := exec.Command("git", "show", commitId, "--decorate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func gitPushTags(remote string, tags ...string) (string, error) {
	args := []string{"push", remote}
	args = append(args, tags...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
