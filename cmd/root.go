package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "veye",
	Short: "A brief description of your application",
	Long:  `The Vuln_EYE Tool is a command-line utility for detecting vulnerabilities in Android apps. It scans Android manifests, Gradle files, and pubspec.yaml files, and identifies vulnerable dependencies. It integrates seamlessly into CI pipelines for automated vulnerability detection.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
