package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Uttkarsh-raj/gitup/helper"
	"github.com/spf13/cobra"
)

var analyze = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repo",
	Long:  `Scan the code to find the presence of any vulnerable dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("./app/build.gradle")
		// file, err := os.Open("test.gradle")
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lines := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		// Generate the map of dependencies
		depMap := helper.GetDependency(lines)

		resp, err := helper.GetData(depMap) // check for errors
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Print(resp) // Currently printing but i dont think we need this, we can make some use of this later

		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyze)
}
