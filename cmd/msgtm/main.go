package main

import (
	"fmt"
	//	"os/exec"

	"github.com/spf13/cobra"
)

type CommitId string

const HEAD CommitId = "HEAD"

func main() {
	services := []string{}
	commitId := new(string)
	tagVersion := new(string)
	isAll := new(bool)
	patch := new(bool)
	minor := new(bool)
	major := new(bool)
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
				println("tag all services")
			}
		},
	}

	tagCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of services")
	tagCmd.Flags().StringVarP(tagVersion, "version", "v", "", "Tag version")
	tagCmd.Flags().BoolVarP(patch, "patch", "p", false, "Patch version up")
	tagCmd.Flags().BoolVarP(minor, "minor", "m", false, "Minor version up")
	tagCmd.Flags().BoolVarP(major, "major", "M", false, "Major version up")
	tagCmd.Flags().StringVarP(commitId, "commit-id", "c", "", "Commit ID")
	tagCmd.Flags().BoolVarP(isAll, "all", "a", false, "Tag all services")

	rootCmd.AddCommand(tagCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
