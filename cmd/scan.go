package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Uttkarsh-raj/veye/helper"
	"github.com/spf13/cobra"
)

var analyze = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repo",
	Long:  `Scan the code to find the presence of any vulnerable dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		isFlutter, _ := cmd.Flags().GetBool("flutter") // default then android else when --flutter is added run for flutter
		verbose, _ := cmd.Flags().GetBool("verbose")
		fix, _ := cmd.Flags().GetBool("fix")
		if isFlutter {
			err := helper.ScanFlutterProject("./pubspec.yaml", verbose, fix)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stdout, "All Dependencies are clean.\n")
				os.Exit(0)
			}
		} else {
			re := regexp.MustCompile(`(?:implementation|testImplementation|androidTestImplementation)\s+(?:platform\()?['"]([^:]+):([^:]+):([^'"]+)['"]\)?`)
			err := helper.ScanProject("./app/build.gradle", re, verbose, fix)
			// err := helper.ScanProject("./test.gradle", re, verbose, fix)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stdout, "All Dependencies are clean.\n")
				os.Exit(0)
			}
		}
	},
}

func init() {
	analyze.Flags().Bool("flutter", false, "Scan for Flutter dependencies")
	analyze.Flags().Bool("verbose", false, "Gives a Verbose output")
	analyze.Flags().Bool("fix", false, "Gives the fixed version if present")
	rootCmd.AddCommand(analyze)
}
