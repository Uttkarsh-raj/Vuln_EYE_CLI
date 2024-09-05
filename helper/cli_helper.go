package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/Uttkarsh-raj/gitup/models"
)

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

// Get data from the dependencies json
func GetData(dependencyMap map[string]string) (string, error) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(dependencyMap)) // trap the errors
	response := ""

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
				errorChannel <- fmt.Errorf("error: Error marshaling payload for %s:%s: %w", name, version, err)
				return
			}

			// Sending POST request
			url := "https://api.osv.dev/v1/query"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				errorChannel <- fmt.Errorf("error: Error creating HTTP request for %s:%s: %w", name, version, err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errorChannel <- fmt.Errorf("error:Error sending HTTP request for %s:%s: %w", name, version, err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errorChannel <- fmt.Errorf("error: Error reading response body for %s:%s: %w", name, version, err)
				return
			}

			response += fmt.Sprintf("Response for %s:%s => %s\n", name, version, string(body)) + "\n"

			if len(string(body)) > 5 {
				errorChannel <- fmt.Errorf("error: Vulnerability found in this package for %s:%s", name, version)
				return
			}

		}(name, version)
	}

	wg.Wait()
	close(errorChannel)

	var finalError string
	for err := range errorChannel {
		if finalError == "" {
			finalError = err.Error()
		} else {
			finalError += "\n" + err.Error()
		}
	}

	if finalError != "" {
		return response, fmt.Errorf(finalError)
	}

	return response, nil
}
