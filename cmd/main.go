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

	// tag add
	tagAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add is a tool for multi service git tag manager",
		Run:   tagAddCmd(),
	}
	tagAddCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagAddCmd.Flags().StringSliceP("services", "s", []string{}, "Add of services")
	tagAddCmd.Flags().StringP("from-config-file", "f", "", "Add of services from config file")

	// tag version-up
	tagVersionUpCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "version-up is a tool for multi service git tag manager",
		Run:   tagVersionUpCmd(),
	}
	tagVersionUpCmd.Flags().BoolP("minor", "m", false, "Minor version up")
	tagVersionUpCmd.Flags().BoolP("major", "M", false, "Major version up")
	tagVersionUpCmd.Flags().BoolP("all", "a", false, "Tag all services")
	tagVersionUpCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagVersionUpCmd.Flags().StringSliceP("services", "s", []string{}, "List of services")

	// tag reset
	tagResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "reset is a tool for multi service git tag manager",
		Run:   tagResetCmd(),
	}
	tagResetCmd.Flags().BoolP("origin", "o", false, "Reset origin")
	tagResetCmd.Flags().BoolP("exclude-local", "e", false, "Exclude local")

	// tags push
	tagsPushCmd := &cobra.Command{
		Use:   "push",
		Short: "push is a tool for multi service git tag manager",
		Run:   tagsPushCmd(),
	}
	tagsPushCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagsPushCmd.Flags().StringP("remote", "r", "", "Remote")

	rootCmd.AddCommand(tagAddCmd)
	rootCmd.AddCommand(tagVersionUpCmd)
	rootCmd.AddCommand(tagResetCmd)
	rootCmd.AddCommand(tagsPushCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

type CobraCmdRunner func(cmd *cobra.Command, args []string)

func tagsPushCmd() CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		commitIdStr, _ := cmd.Flags().GetString("commit-id")
		remoteStr, _ := cmd.Flags().GetString("remote")

		commitId := msgtm.HEAD
		if commitIdStr != "" {
			commitId = msgtm.CommitId(commitIdStr)
		}
		remote := msgtm.Origin
		if remoteStr != "" {
			remote = msgtm.RemoteAddr(remoteStr)
		}

		getter := msgtm.DefaultCommitTagGetter()
		pusher := msgtm.DefaultGitTagPusher()

		err := msgtm.PushAll(
			getter,
			pusher,
			&remote,
			&commitId,
		)
		if err != nil {
			fmt.Printf("Failed to push service tags: %s\n", err.Error())
			return
		}
	}
}

func tagResetCmd() CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		commitId := msgtm.HEAD
		if len(args) > 0 {
			commitId = msgtm.CommitId(args[0])
		}

		getter := msgtm.DefaultCommitTagGetter()
		destroyer := &msgtm.DestroyDecorator{}
		origin, _ := cmd.Flags().GetBool("origin")
		if origin {
			destroyer.Clients = append(destroyer.Clients, msgtm.ForceOriginDestroyer())
		}
		excludeLocal, _ := cmd.Flags().GetBool("exclude-local")
		if !excludeLocal {
			destroyer.Clients = append(destroyer.Clients, msgtm.ForceDestroyer())
		}

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
