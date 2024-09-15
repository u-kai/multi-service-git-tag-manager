package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

type CommitId string

const HEAD CommitId = "HEAD"

type SemVer struct {
	major int
	minor int
	patch int
}

func FromStr(s string) (SemVer, error) {
	var major, minor, patch int
	_, err := fmt.Sscanf(s, "v%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		_, err = fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return SemVer{}, fmt.Errorf("invalid semver string: %s", s)
		}
	}
	return SemVer{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}

func (s SemVer) MajorUp() SemVer {
	return SemVer{
		major: s.major + 1,
		minor: 0,
		patch: 0,
	}
}

func (s SemVer) MinorUp() SemVer {
	return SemVer{
		major: s.major,
		minor: s.minor + 1,
		patch: 0,
	}
}

func (s SemVer) PatchUp() SemVer {
	return SemVer{
		major: s.major,
		minor: s.minor,
		patch: s.patch + 1,
	}
}

func (s *SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

type ServiceTagWithSemVer struct {
	service string
	tag     SemVer
}

func NewServiceTagWithSemVer(service string, tag SemVer) *ServiceTagWithSemVer {
	return &ServiceTagWithSemVer{
		service: service,
		tag:     tag,
	}
}

func (s *ServiceTagWithSemVer) String() string {
	return fmt.Sprintf("%s-%s", s.service, s.tag.String())
}

func main() {
	services := []string{}
	commitId := new(string)
	tagVersion := new(string)
	isAll := new(bool)
	rootCmd := &cobra.Command{
		Use:   "msgtm",
		Short: "msgtm is a tool for multi service git tag manager",
	}
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "tag is a tool for multi service git tag manager",
		Run: func(cmd *cobra.Command, args []string) {
			if *isAll {
				fmt.Println("tag all services")
			} else {
				if *commitId == "" {
					*commitId = "HEAD"
				}
				for _, service := range services {
					tag := *tagVersion
					semVer, err := FromStr(tag)
					if err == nil {
						tag = semVer.String()
					}
					serviceTag := NewServiceTagWithSemVer(service, semVer)

					message := fmt.Sprintf(`"create auto tag:%s-%s"`, service, tag)

					gitTagCmd := exec.Command("git", "tag", "-a", serviceTag.String(), *commitId, "-m", message)
					c := gitTagCmd.String()
					println(c)
					output, err := gitTagCmd.CombinedOutput()
					if err != nil {
						println(err.Error())
					}
					println(output)
					fmt.Println(string(output))
				}
			}
		},
	}

	tagCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of services")
	tagCmd.Flags().StringVarP(tagVersion, "version", "v", "", "Tag version")
	commitId = tagCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	isAll = tagCmd.Flags().BoolP("all", "a", false, "Tag all services")

	rootCmd.AddCommand(tagCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
