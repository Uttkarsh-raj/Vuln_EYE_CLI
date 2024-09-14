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

// Function to generate a map of dependencies
func GetDependency(lines []string, re *regexp.Regexp) map[string]string {
	deps := make(map[string]string)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 3 {
			// Construct key as "group:artifact"
			key := matches[1] + ":" + matches[2]
			version := matches[3]
			deps[key] = version // key is package name, value is version
		}
	}
	return deps
}

// Get data from the dependencies json
func GetData(dependencyMap map[string]string, verbose bool, fix bool) (string, error) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(dependencyMap)) // trap the errors
	report := ""                                         // To store the final report

	for name, version := range dependencyMap {
		wg.Add(1)

		go func(name, version string) {
			defer wg.Done()

			// Creating the payload
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
				errorChannel <- fmt.Errorf("error: Error sending HTTP request for %s:%s: %w", name, version, err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errorChannel <- fmt.Errorf("error: Error reading response body for %s:%s: %w", name, version, err)
				return
			}

			// Unmarshal response to check for vulnerabilities
			var response models.VerboseResp
			err = json.Unmarshal(body, &response)
			if err != nil {
				errorChannel <- fmt.Errorf("error: Error unmarshalling response body: %s", err.Error())
				return
			}

			// Generate report for each dependency
			if len(response.Vulns) == 0 {
				// Clean dependency (Green)
				report += fmt.Sprintf("\033[32m✅ %s (Version: %s) is clean.\033[0m\n", name, version)
			} else {
				// Vulnerable dependency (Red)
				vulnerableReport := fmt.Sprintf("\033[31m❌ %s (Version: %s) has vulnerabilities:\n", name, version)

				// Verbose flag shows ID and Summary
				if verbose {
					for _, vuln := range response.Vulns {
						vulnerableReport += fmt.Sprintf("  - ID: %s\n  - Summary: %s\n", vuln.Aliases[0], vuln.Summary)
					}
				}

				// Fix flag shows Introduced and Fixed versions
				if fix {
					for _, vuln := range response.Vulns {
						for _, affected := range vuln.Affected {
							for _, r := range affected.Ranges {
								for _, event := range r.Events {
									if event.Introduced != "" {
										vulnerableReport += fmt.Sprintf("  - Introduced: %s", event.Introduced)
									}
									if event.Fixed != "" {
										vulnerableReport += fmt.Sprintf("\n  - Fixed: %s\n", event.Fixed)
									}
								}
							}
						}
					}
				}

				vulnerableReport += "\033[0m" // Reset color at the end
				errorChannel <- fmt.Errorf(vulnerableReport)
				report += vulnerableReport // Append to final report
			}
		}(name, version)
	}

	wg.Wait()
	close(errorChannel)

	// Collect errors to determine if CI should fail
	var finalError string
	for err := range errorChannel {
		if finalError == "" {
			finalError = err.Error()
		} else {
			finalError += "\n" + err.Error()
		}
	}

	// Print the final report
	fmt.Println("\n")
	fmt.Println("===== Dependency Scan Report =====")
	fmt.Println(report)
	fmt.Println("===== End of Report =====")

	if finalError != "" {
		return report, fmt.Errorf(finalError) // Fails CI with errors
	}

	return report, nil // CI passes with a clean report
}
