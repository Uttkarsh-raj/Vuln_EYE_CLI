package helper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/Uttkarsh-raj/gitup/models"
	"gopkg.in/yaml.v3"
)

// Function to generate map of dependencies
func GetDependency(lines []string, re *regexp.Regexp) map[string]string {

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

func ScanProject(fileLoc string, re *regexp.Regexp) error {
	file, err := os.Open(fileLoc)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Generate the map of dependencies
	depMap := GetDependency(lines, re)

	resp, err := GetData(depMap) // check for errors
	if err != nil {
		return err
	}

	fmt.Print(resp) // Currently printing but i dont think we need this, we can make some use of this later

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func ScanFlutterProject(fileLoc string) error {
	file, err := os.Open(fileLoc)
	if err != nil {
		return err
	}
	defer file.Close()

	var pubspec models.Pubspec
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&pubspec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Convert dependencies
	deps := ConvertVersions(pubspec.Dependencies)
	resp, err := GetData(deps) // check for errors
	if err != nil {
		return err
	}
	fmt.Println("Dependencies: ")
	fmt.Print(resp)

	// Convert dev dependencies
	devDeps := ConvertVersions(pubspec.DevDependencies)
	resp, err = GetData(devDeps) // check for errors
	if err != nil {
		return err
	}
	fmt.Println("Dev Dependencies: ")
	fmt.Print(resp)

	return nil
}

func ConvertVersions(deps map[string]interface{}) map[string]string {
	normalized := make(map[string]string)
	for pkg, ver := range deps {
		var version string
		switch v := ver.(type) {
		case string:
			version = normalizeVersion(v)
		case map[interface{}]interface{}:
			if len(v) == 1 {
				for _, value := range v {
					if strValue, ok := value.(string); ok {
						version = normalizeVersion(strValue)
					}
				}
			}
		}
		if version != "" {
			normalized[pkg] = version
		}
	}
	return normalized
}

// normalizeVersion extracts and normalizes the version string from various formats.
func normalizeVersion(version string) string {
	version = strings.TrimLeft(version, "^>=~")
	re := regexp.MustCompile(`^(\d+\.\d+\.\d+)`)
	match := re.FindStringSubmatch(version)
	if len(match) > 0 {
		return match[1]
	}
	return version
}
