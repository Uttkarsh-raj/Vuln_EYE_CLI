package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/Uttkarsh-raj/gitup/models"
	"github.com/spf13/cobra"
)

var analyze = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repo",
	Long:  `Scan the code to find the presence of any vulnerable dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("test.txt")
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
		depMap := GetDependency(lines)

		getData(depMap)

		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyze)
}

// Function to generate map of dependencies
func GetDependency(lines []string) map[string]string {
	re := regexp.MustCompile(`(?:implementation|testImplementation|androidTestImplementation)\s+(?:platform\()?['"]([^:]+):([^:]+):([^'"]+)['"]\)?`)

	deps := make(map[string]string)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 3 {
			// Construct key as "group:artifact"
			key := matches[1] + ":" + matches[2]
			version := matches[3]
			deps[key] = version //key is package name and value is version
		}
	}
	return deps
}

func getData(dependencyMap map[string]string) {
	var wg sync.WaitGroup

	for name, version := range dependencyMap {
		wg.Add(1)

		go func(name, version string) {
			defer wg.Done()

			// Creating the paylaod
			payload := models.DepsPayload{
				Version: version,
			}
			payload.Package.Name = name

			// Marshaling payload into JSON
			jsonData, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error marshaling payload:", err)
				return
			}

			// Sending POST request
			url := "https://api.osv.dev/v1/query"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating HTTP request:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending HTTP request:", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			fmt.Printf("Response for %s:%s => %s\n", name, version, string(body))
		}(name, version)
	}

	wg.Wait()
}
