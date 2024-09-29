package executor

import (
	"log/slog"
	"msgtm/pkg/domain"
	"os/exec"
	"strings"
)

type gitCommandExecutor func(args ...string) (string, error)

func GitShellCommandExecutor() gitCommandExecutor {
	return func(args ...string) (string, error) {
		cmd := exec.Command("git", args...)
		output, err := cmd.CombinedOutput()
		return string(output), err
	}
}

func LogDecorateToExecutor(gitCmd gitCommandExecutor, logger slog.Logger) gitCommandExecutor {
	return func(args ...string) (string, error) {
		logger.Debug("git command", slog.Any("args", args))
		output, err := gitCmd(args...)
		if err != nil {
			logger.Error("git command failed", slog.Any("error", err), slog.String("output", output))
			return output, err
		}
		logger.Debug("git command", slog.String("output", output))
		return output, nil
	}
}

func tagList(executor gitCommandExecutor) (*[]domain.GitTag, error) {
	output, err := executor("tag")
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

func gitTagAddLight(executor gitCommandExecutor, commitId string, tag string) (string, error) {
	return executor("tag", tag, commitId)
}

func gitTagAdd(executor gitCommandExecutor, commitId string, tag string, message string) (string, error) {
	return executor("tag", "-a", tag, "-m", message, commitId)
}

func gitTagDelete(executor gitCommandExecutor, tag string, force bool) (string, error) {
	deleteOption := "-d"
	if force {
		deleteOption = "-d"
	}
	return executor("tag", deleteOption, tag)
}

func gitTagRemoteDelete(executor gitCommandExecutor, remote string, tags []string) (string, error) {
	deleteOption := "--delete"
	cmdArgs := []string{"push", remote, deleteOption}
	cmdArgs = append(cmdArgs, tags...)
	return executor(cmdArgs...)
}

func gitShowCommitTags(executor gitCommandExecutor, commitId string) ([]string, error) {
	output, err := gitShowCommit(executor, commitId)
	if err != nil {
		return nil, err
	}
	commitLine := strings.Split(output, "\n")[0]
	// tag: service1-v1.1.1, tags: service2-v1.1.1
	tagsStr := strings.Split(commitLine, "(")[1]
	// remove ")"
	tagsStr = tagsStr[:len(tagsStr)-1]
	result := []string{}
	for _, tagStr := range strings.Split(tagsStr, ", ") {
		if !strings.HasPrefix(tagStr, "tag: ") {
			continue
		}
		result = append(result, strings.Split(tagStr, "tag: ")[1])
	}
	return result, nil
}

func gitShowCommit(executor gitCommandExecutor, commitId string) (string, error) {
	return executor("show", commitId, "--decorate")
}

func gitPushTags(executor gitCommandExecutor, remote string, tags ...string) (string, error) {
	args := []string{"push", remote}
	args = append(args, tags...)
	return executor(args...)
}
