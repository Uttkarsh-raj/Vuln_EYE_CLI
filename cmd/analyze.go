package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Uttkarsh-raj/gitup/helper"
	"github.com/spf13/cobra"
)

var analyze = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repo",
	Long:  `Scan the code to find the presence of any vulnerable dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		isFlutter, _ := cmd.Flags().GetBool("flutter") // default then android else when --flutter is added run for flutter
		if isFlutter {
			err := helper.ScanFlutterProject("./pubspec.yaml")
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
		} else {
			re := regexp.MustCompile(`(?:implementation|testImplementation|androidTestImplementation)\s+(?:platform\()?['"]([^:]+):([^:]+):([^'"]+)['"]\)?`)
			err := helper.ScanProject("./app/build.gradle", re)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	analyze.Flags().Bool("flutter", false, "Scan for Flutter dependencies")
	rootCmd.AddCommand(analyze)
}
