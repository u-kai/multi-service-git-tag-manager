package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

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
					message := fmt.Sprintf(`"create auto tag:%s-%s"`, service, *tagVersion)

					gitTagCmd := exec.Command("git", "tag", "-a", service+"-"+*tagVersion, *commitId, "-m", message)
					c := gitTagCmd.String()
					println(c)
					output, err := gitTagCmd.CombinedOutput()
					if err != nil {
						println("error")
						println(string(output))
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
