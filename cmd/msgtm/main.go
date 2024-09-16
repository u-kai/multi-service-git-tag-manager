package main

import (
	"fmt"
	"msgtm"
	"os"

	//	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type CommitId string

const HEAD CommitId = "HEAD"

func main() {
	rootCmd := &cobra.Command{
		Use:   "msgtm",
		Short: "msgtm is a tool for multi service git tag manager",
	}

	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "tag is a tool for multi service git tag manager",
	}

	// tag add
	tagAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add is a tool for multi service git tag manager",
		Run:   tagAddCmd(),
	}
	tagAddCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagAddCmd.Flags().StringSliceP("services", "s", []string{}, "Add of services")
	tagAddCmd.Flags().StringP("from-config-file", "f", "", "Add of services from config file")
	tagCmd.AddCommand(tagAddCmd)

	tagVersionUpCmd := &cobra.Command{
		Use:   "up",
		Short: "version-up is a tool for multi service git tag manager",
		Run:   tagVersionUpCmd(),
	}

	tagVersionUpCmd.Flags().BoolP("minor", "m", false, "Minor version up")
	tagVersionUpCmd.Flags().BoolP("major", "M", false, "Major version up")
	tagVersionUpCmd.Flags().BoolP("all", "a", false, "Tag all services")
	tagVersionUpCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagVersionUpCmd.Flags().StringSliceP("services", "s", []string{}, "List of services")
	tagCmd.AddCommand(tagVersionUpCmd)

	tagResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "reset is a tool for multi service git tag manager",
		Run:   tagResetCmd(),
	}
	tagCmd.AddCommand(tagResetCmd)

	rootCmd.AddCommand(tagCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

type CobraCmdRunner func(cmd *cobra.Command, args []string)

func tagResetCmd() CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		commitId := msgtm.HEAD
		if len(args) > 0 {
			commitId = msgtm.CommitId(args[0])
		}

		getter := msgtm.DefaultCommitTagGetter()
		destroyer := msgtm.ForceDestroyer()

		err := msgtm.ResetServiceTags(
			destroyer,
			getter,
			&commitId,
		)
		if err != nil {
			fmt.Printf("Failed to reset service tags: %s\n", err.Error())
			return
		}
	}
}

func tagVersionUpCmd() CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		minor, _ := cmd.Flags().GetBool("minor")
		major, _ := cmd.Flags().GetBool("major")
		isAll, _ := cmd.Flags().GetBool("all")
		commitIdStr, _ := cmd.Flags().GetString("commit-id")
		services, _ := cmd.Flags().GetStringSlice("services")

		commitId := msgtm.HEAD
		if commitIdStr != "" {
			commitId = msgtm.CommitId(commitIdStr)
		}

		var list msgtm.TagList = &msgtm.AllTagList{}
		if !isAll && len(services) > 0 {
			list = &msgtm.FilterTagList{
				IncludePrefix: services,
			}
		}

		f := msgtm.PatchUpAll
		if minor {
			f = msgtm.MinorUpAll
		}
		if major {
			f = msgtm.MajorUpAll
		}

		register := msgtm.DefaultGitTagRegister()

		err := msgtm.VersionUpAllServiceTags(
			list,
			register,
			f,
			&commitId,
		)

		if err != nil {
			fmt.Printf("Failed to version up all service tags: %s\n", err.Error())
			return
		}
	}
}

func tagAddCmd() CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: tag add command must version args.")
			return
		}
		version := args[0]
		semVer, err := msgtm.FromStr(version)
		if err != nil {
			fmt.Printf("Failed to parse args: %s\n err msg: %s", version, err.Error())
			return
		}
		commitIdStr, _ := cmd.Flags().GetString("commit-id")
		services, _ := cmd.Flags().GetStringSlice("services")
		fileName, _ := cmd.Flags().GetString("from-config-file")

		if len(fileName) > 0 {
			content, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Printf("Failed to read file: %s\n", err.Error())
				return
			}
			config := ServiceConfig{}
			err = yaml.Unmarshal(content, &config)
			if err != nil {
				fmt.Printf("Failed to unmarshal yaml: %s\n", err.Error())
				return
			}
			for _, service := range config.Services {
				services = append(services, service.Name)
			}
		}

		commitId := msgtm.HEAD
		if commitIdStr != "" {
			commitId = msgtm.CommitId(commitIdStr)
		}

		register := msgtm.DefaultGitTagRegister()
		err = msgtm.CreateServiceTags(
			register,
			&commitId,
			services,
			semVer,
		)
		if err != nil {
			fmt.Printf("Failed to create service tags: %s\n", err.Error())
			return
		}
	}
}

type ServiceConfig struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
}
