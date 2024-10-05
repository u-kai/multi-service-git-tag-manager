package main

import (
	"fmt"
	"io"
	"log/slog"
	"msgtm/pkg/domain"
	"msgtm/pkg/executor"
	"msgtm/pkg/subcmd"
	"msgtm/pkg/usecase"
	"os"
	"strings"

	"github.com/spf13/cobra"
	//"gopkg.in/yaml.v2"
)

type CommitId string

const HEAD CommitId = "HEAD"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	gitExecutor := executor.LogDecorateToExecutor(
		executor.GitShellCommandExecutor(),
		logger,
		func(output string) string {
			split := strings.Split(output, "\n")
			return split[0] + " ... " + "output line length: " + fmt.Sprintf("%d", len(split))
		},
	)

	getter := &executor.LoggingQueryExecutor[usecase.GetCommitTagQuery, *[]domain.GitTag]{
		Executor: &executor.CommitTagGetter{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}
	register := &executor.LoggingCommandExecutor[usecase.RegisterServiceTagsCommand]{
		Executor: executor.NewGitTagRegister(gitExecutor),
		Logger:   logger,
	}
	list := &executor.LoggingQueryExecutor[usecase.ListTagsQuery, *[]domain.GitTag]{
		Executor: &executor.GitTagList{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}
	localDestroyer := &executor.LoggingCommandExecutor[usecase.DestroyServiceTagsCommand]{
		Executor: &executor.LocalServiceTagsDestroyer{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}
	remoteDestroyer := &executor.LoggingCommandExecutor[usecase.DestroyServiceTagsCommand]{
		Executor: &executor.RemoteServiceTagsDestroyer{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}
	pusher := &executor.LoggingCommandExecutor[usecase.CommitPushCommand]{
		Executor: &executor.GitTagPusher{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}
	finder := &executor.LoggingQueryExecutor[usecase.FindCommitQuery, *domain.CommitId]{
		Executor: &executor.CommitFinder{
			GitCommandExecutor: gitExecutor,
		},
		Logger: logger,
	}

	rootCmd := &cobra.Command{
		Use:   "msgtn",
		Short: "msgtn is a tool for multi service git tag manager",
	}

	rootCmd.AddCommand(listCmd(logger, list, finder))
	rootCmd.AddCommand(tagAddCmd(logger, register, list, finder))
	rootCmd.AddCommand(tagVersionUpCmd(logger, list, register, getter, finder))
	rootCmd.AddCommand(tagResetCmd(logger, getter, localDestroyer, remoteDestroyer, list, finder))
	rootCmd.AddCommand(tagsPushCmd(logger, getter, pusher))
	rootCmd.AddCommand(syncAllCmd(list, finder))
	rootCmd.AddCommand(initCmd(logger))

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

type CobraCmdRunner func(cmd *cobra.Command, args []string)

func initCmd(logger *slog.Logger) *cobra.Command {
	f := func() CobraCmdRunner {
		return func(cmd *cobra.Command, args []string) {
			fileName, _ := cmd.Flags().GetString("filename")
			services, _ := cmd.Flags().GetStringSlice("services")
			serviceConfigs := make([]domain.ServiceName, 0)
			for _, service := range services {
				serviceConfigs = append(serviceConfigs, domain.ServiceName(service))
			}
			logger.Debug("serviceConfigs", slog.Any("serviceConfigs", serviceConfigs))
			stateWriter := domain.InitStateWriter(serviceConfigs...)
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Failed to create file: %s\n", err.Error())
				return
			}
			err = stateWriter.Write(file, domain.YAML)
			if err != nil {
				fmt.Printf("Failed to write file: %s\n", err.Error())
				return
			}
		}
	}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init is a tool for multi service git tag manager",
		Run:   f(),
	}
	initCmd.Flags().StringP("filename", "f", "services-state.yaml", "filename")
	initCmd.Flags().StringSliceP("services", "s", []string{}, "services")
	return initCmd
}

func syncAll(writer io.Writer, state *domain.WritedState, list usecase.ListTags, finder usecase.CommitFinder) error {
	state, err := usecase.SyncAllServiceTagState(state, list, finder)
	if err != nil {
		return err
	}
	err = state.Write(writer, domain.YAML)
	if err != nil {
		return err
	}
	return nil
}

func addSyncAll(
	f CobraCmdRunner,
	list usecase.ListTags,
	finder usecase.CommitFinder,
) CobraCmdRunner {
	return func(cmd *cobra.Command, args []string) {
		f(cmd, args)
		sync, _ := cmd.Flags().GetBool("sync")
		fileName, _ := cmd.Flags().GetString("state-file")
		if sync {
			file, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Printf("Failed to open file: %s\n", err.Error())
				return
			}
			state, err := domain.FromReader(file, domain.YAML)
			if err != nil {
				fmt.Printf("Failed to read file: %s\n", err.Error())
				return
			}
			err = syncAll(file, state, list, finder)
			if err != nil {
				fmt.Printf("Failed to sync all service tags: %s\n", err.Error())
				return
			}
		}
	}
}

func syncAllCmd(list usecase.ListTags, finder usecase.CommitFinder) *cobra.Command {
	f := addSyncAll(func(_ *cobra.Command, _ []string) {}, list, finder)
	syncAllCmd := &cobra.Command{
		Use:   "sync",
		Short: "sync is a tool for multi service git tag manager",
		Run:   f,
	}
	syncAllCmd.Flags().Bool("sync", true, "Sync all service tags")
	syncAllCmd.Flags().StringP("state-file", "t", "services-state.yaml", "State file")
	return syncAllCmd
}

func listCmd(logger *slog.Logger, list usecase.ListTags, finder usecase.CommitFinder) *cobra.Command {
	f := func(list usecase.ListTags, finder usecase.CommitFinder) CobraCmdRunner {
		return func(cmd *cobra.Command, args []string) {
			services, _ := cmd.Flags().GetStringSlice("services")
			isAll, _ := cmd.Flags().GetBool("isAll")
			err := subcmd.LogSubCommandDecorator(
				subcmd.ServiceTagsListCommand(list, finder),
				logger,
			)(subcmd.ServiceTagsListParameter{
				Filter: services,
				IsAll:  isAll,
			})
			if err != nil {
				fmt.Printf("Failed to list service tags: %s\n", err.Error())
				return
			}
		}
	}
	serviceTagsListCmd := &cobra.Command{
		Use:   "list",
		Short: "list is a tool for multi service git tag manager",
		Run:   f(list, finder),
	}
	serviceTagsListCmd.Flags().StringSliceP("services", "s", []string{}, "services")
	serviceTagsListCmd.Flags().Bool("isAll", true, "List all service tags")
	return serviceTagsListCmd
}

func tagAddCmd(logger *slog.Logger, register usecase.RegisterServiceTags, list usecase.ListTags, finder usecase.CommitFinder) *cobra.Command {
	f := func(register usecase.RegisterServiceTags) CobraCmdRunner {
		return func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Error: tag add command must version args.")
				return
			}
			version := args[0]
			commitIdStr, _ := cmd.Flags().GetString("commit-id")
			services, _ := cmd.Flags().GetStringSlice("services")
			fileName, _ := cmd.Flags().GetString("from-config-file")

			param := subcmd.TagAddCommandParameter{
				Version:        version,
				CommitId:       commitIdStr,
				Services:       services,
				FromConfigFile: fileName,
			}

			err := subcmd.LogSubCommandDecorator(
				subcmd.TagAddCommand(register),
				logger,
			)(param)

			if err != nil {
				fmt.Printf("Failed to add service tags: %s\n", err.Error())
			}

		}
	}
	tagAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add is a tool for multi service git tag manager",
		Run:   addSyncAll(f(register), list, finder),
	}
	tagAddCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagAddCmd.Flags().StringSliceP("services", "s", []string{}, "Add of services")
	tagAddCmd.Flags().StringP("from-config-file", "f", "", "Add of services from config file")
	tagAddCmd.Flags().Bool("sync", true, "Sync all service tags")
	tagAddCmd.Flags().StringP("state-file", "t", "services-state.yaml", "State file")
	return tagAddCmd
}
func tagsPushCmd(logger *slog.Logger, getter usecase.CommitTagGetter, pusher usecase.CommitPusher) *cobra.Command {
	f := func(getter usecase.CommitTagGetter, pusher usecase.CommitPusher) CobraCmdRunner {
		return func(cmd *cobra.Command, args []string) {
			commitIdStr, _ := cmd.Flags().GetString("commit-id")
			remoteStr, _ := cmd.Flags().GetString("remote")

			param := subcmd.PushCommandParameter{
				CommitId: commitIdStr,
				Remote:   remoteStr,
			}
			err := subcmd.LogSubCommandDecorator(
				subcmd.PushCommand(getter, pusher),
				logger,
			)(param)
			if err != nil {
				fmt.Printf("Failed to push service tags: %s\n", err.Error())
				return
			}
		}
	}
	tagsPushCmd := &cobra.Command{
		Use:   "push",
		Short: "push is a tool for multi service git tag manager",
		Run:   f(getter, pusher),
	}
	tagsPushCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagsPushCmd.Flags().StringP("remote", "r", "", "Remote")
	return tagsPushCmd
}

func tagResetCmd(logger *slog.Logger, getter usecase.CommitTagGetter, localDestroyer usecase.DestroyServiceTags, remoteDestroyer usecase.DestroyServiceTags, list usecase.ListTags, finder usecase.CommitFinder) *cobra.Command {
	f := func(cmd *cobra.Command, args []string) {
		origin, _ := cmd.Flags().GetBool("origin")
		excludeLocal, _ := cmd.Flags().GetBool("exclude-local")
		commitIdStr, _ := cmd.Flags().GetString("commit-id")
		param := subcmd.ResetCommandParameter{
			Origin:       origin,
			ExcludeLocal: excludeLocal,
			CommitId:     commitIdStr,
		}
		if len(args) > 0 {
			param.CommitId = args[0]
		}

		err := subcmd.LogSubCommandDecorator(
			subcmd.ResetCommand(getter, localDestroyer, remoteDestroyer),
			logger,
		)(param)
		if err != nil {
			fmt.Printf("Failed to reset service tags: %s\n", err.Error())
		}
	}
	tagResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "reset is a tool for multi service git tag manager",
		Run:   addSyncAll(f, list, finder),
	}
	tagResetCmd.Flags().BoolP("origin", "o", false, "Reset origin")
	tagResetCmd.Flags().BoolP("exclude-local", "e", false, "Exclude local")
	tagResetCmd.Flags().StringP("state-file", "f", "services-state.yaml", "State file")
	tagResetCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagResetCmd.Flags().Bool("sync", true, "Sync all service tags")
	return tagResetCmd
}

func tagVersionUpCmd(logger *slog.Logger, list usecase.ListTags, register usecase.RegisterServiceTags, getter usecase.CommitTagGetter, finder usecase.CommitFinder) *cobra.Command {
	f := func(cmd *cobra.Command, args []string) {
		minor, _ := cmd.Flags().GetBool("minor")
		major, _ := cmd.Flags().GetBool("major")
		isAll, _ := cmd.Flags().GetBool("all")
		commitIdStr, _ := cmd.Flags().GetString("commit-id")
		services, _ := cmd.Flags().GetStringSlice("services")

		param := subcmd.VersionUpCommandParameter{
			Minor:    minor,
			Major:    major,
			IsAll:    isAll,
			CommitId: commitIdStr,
			Services: services,
		}

		err := subcmd.LogSubCommandDecorator(
			subcmd.VersionUpCommand(
				list,
				register,
				getter,
			),
			logger,
		)(param)

		if err != nil {
			fmt.Printf("Failed to version up all service tags: %s\n", err.Error())
			return
		}
	}
	tagVersionUpCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "version-up is a tool for multi service git tag manager",
		Run:   addSyncAll(f, list, finder),
	}
	tagVersionUpCmd.Flags().BoolP("minor", "m", false, "Minor version up")
	tagVersionUpCmd.Flags().BoolP("major", "M", false, "Major version up")
	tagVersionUpCmd.Flags().BoolP("all", "a", false, "Tag all services")
	tagVersionUpCmd.Flags().StringP("commit-id", "c", "", "Commit ID")
	tagVersionUpCmd.Flags().StringSliceP("services", "s", []string{}, "List of services")
	tagVersionUpCmd.Flags().Bool("sync", true, "Sync all service tags")
	tagVersionUpCmd.Flags().StringP("state-file", "t", "services-state.yaml", "State file")
	return tagVersionUpCmd
}

type ServiceConfig struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
}
