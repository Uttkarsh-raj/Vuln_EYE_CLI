package helper

import (
	"bufio"
	"log"
	"os"
	"regexp"

	"github.com/Uttkarsh-raj/gitup/models"
	"gopkg.in/yaml.v3"
)

// Function for scanning Android Projects.
func ScanProject(fileLoc string, re *regexp.Regexp, verbose bool, fix bool) error {
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

	if err := AnalyzeManifest(); err != nil {
		return err
	}

	_, err = GetData(depMap, verbose, fix) // check for errors
	if err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Function for scanning Flutter Projects
func ScanFlutterProject(fileLoc string, verbose bool, fix bool) error {
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
	_, err = GetData(deps, verbose, fix) // check for errors
	if err != nil {
		return err
	}

	// Convert dev dependencies
	devDeps := ConvertVersions(pubspec.DevDependencies)
	_, err = GetData(devDeps, verbose, fix) // check for errors
	if err != nil {
		return err
	}

	if err := AnalyzeManifest(); err != nil {
		return err
	}

	return nil
}
