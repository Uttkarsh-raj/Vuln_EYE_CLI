package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var dependecies []string

var analyze = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repo",
	Long:  `Scan the code to find the presence of any vulnerable dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("test")
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			dependecy := GetDependency(line)
			if len(dependecy) > 1 {
				dependecies = append(dependecies, dependecy)
			}
		}

		fmt.Println(dependecies) // prints the dependencies

		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyze)
}

func GetDependency(line string) string {

	re := regexp.MustCompile(`(?:implementation|testImplementation|androidTestImplementation)\s+(?:platform\()?['"]([^'"]+)['"]\)?`)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		dependency := matches[1]
		return dependency
	}
	return ""
}
