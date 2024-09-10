package helper

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Uttkarsh-raj/gitup/models"
)

func AnalyzeManifest() error {
	file, err := os.Open("./app/src/main/AndroidManifest.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	manifest := &models.ManifestData{}

	exportedActivityPattern := regexp.MustCompile(`android:exported="(true|false)"`)
	permissionPattern := regexp.MustCompile(`android\.permission\.READ_MEDIA_IMAGES`)
	debuggablePattern := regexp.MustCompile(`android:debuggable="(true|false)"`)
	allowBackupPattern := regexp.MustCompile(`android:allowBackup="(true|false)"`)
	intentFilterPattern := regexp.MustCompile(`<intent-filter>`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "android:exported") {
			matches := exportedActivityPattern.FindStringSubmatch(line)
			if len(matches) == 2 {
				manifest.ExportedActivityFound = true
				if matches[1] == "true" {
					manifest.ExportedActivityTrue = true
				}
			}
		}

		if strings.Contains(line, "<uses-permission") {
			matches := permissionPattern.FindStringSubmatch(line)
			if len(matches) == 2 {
				permission := matches[1]
				manifest.Permissions = append(manifest.Permissions, permission)
			}
		}

		if strings.Contains(line, "android:debuggable") {
			matches := debuggablePattern.FindStringSubmatch(line)
			if len(matches) == 2 {
				manifest.DebuggableFlag = matches[1]
			}
		}

		if strings.Contains(line, "android:allowBackup") {
			matches := allowBackupPattern.FindStringSubmatch(line)
			if len(matches) == 2 {
				manifest.AllowBackupFlag = matches[1]
			}
		}

		if intentFilterPattern.MatchString(line) {
			manifest.IntentFilters = append(manifest.IntentFilters, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	logVulnerabilities(manifest)

	return nil
}

func logVulnerabilities(manifest *models.ManifestData) {
	fmt.Println("===== Manifest Vulnerability Report =====")

	if manifest.ExportedActivityFound && manifest.ExportedActivityTrue {
		fmt.Println("Vulnerability: Exported Activity (android:exported=\"true\")")
		fmt.Println("Severity: High")
		fmt.Println("Risk: The MainActivity is marked as exported, meaning it is accessible to other applications.")
		fmt.Println("       If this activity contains sensitive operations or data, malicious apps could invoke it and exploit the exposed functionality.")
		fmt.Println("Mitigation: Set android:exported=\"false\" if the activity doesn't need to be accessed by other apps.")
		fmt.Println("            If it does need to be exposed, ensure proper intent validation and permission checks.")
		fmt.Println("--------------------------------------------")
	} else if manifest.ExportedActivityFound && !manifest.ExportedActivityTrue {
		fmt.Println("Exported Activity found, but not marked as exported (android:exported=\"false\"). No issue.")
	}

	if len(manifest.Permissions) > 0 {
		fmt.Println("Permissions found in the manifest:")
		for _, permission := range manifest.Permissions {
			fmt.Printf(" - %s\n", permission)
			fmt.Println("  Severity: Unknown (Review if the permission is necessary and justified).")
			fmt.Println("--------------------------------------------")
		}
	} else {
		fmt.Println("No risky permissions found.")
		fmt.Println("--------------------------------------------")
	}

	if manifest.DebuggableFlag == "true" {
		fmt.Println("Vulnerability: android:debuggable=\"true\"")
		fmt.Println("Severity: Medium to High (depending on build)")
		fmt.Println("Risk: If android:debuggable is set to true in production builds, it allows attackers to access the app's debugging information,")
		fmt.Println("       which can be used to inspect app internals or tamper with its behavior.")
		fmt.Println("Mitigation: Ensure android:debuggable=\"false\" is set for production builds.")
		fmt.Println("--------------------------------------------")
	} else if manifest.DebuggableFlag == "false" {
		fmt.Println("No issue: Debuggable flag is set to false.")
	} else {
		fmt.Println("No debuggable flag found.")
	}

	if manifest.AllowBackupFlag == "true" {
		fmt.Println("Vulnerability: Backup Settings (android:allowBackup=\"true\")")
		fmt.Println("Severity: Medium")
		fmt.Println("Risk: The allowBackup flag is set to true, meaning user data and application data can be backed up to cloud services.")
		fmt.Println("       This could expose sensitive data if not encrypted.")
		fmt.Println("Mitigation: Set android:allowBackup=\"false\" if the app handles sensitive data to prevent automatic backups from exposing user information.")
		fmt.Println("--------------------------------------------")
	} else if manifest.AllowBackupFlag == "false" {
		fmt.Println("No issue: Backup is disabled (android:allowBackup=\"false\").")
	} else {
		fmt.Println("No allowBackup flag found.")
	}

	// if len(manifest.IntentFilters) > 0 {
	// 	fmt.Println("Vulnerability: Intent Filters (URL Handling and Sharing Files)")
	// 	fmt.Println("Severity: High")
	// 	fmt.Println("Risk: The app allows sharing of text, images, videos, and generic files through multiple intent-filters.")
	// 	fmt.Println("       Malicious apps could exploit these to send invalid, malicious, or dangerous data.")
	// 	fmt.Println("Mitigation: Validate all input data received through intents to ensure it's safe and follows the expected format. Always sanitize inputs.")
	// 	fmt.Println("--------------------------------------------")
	// } else {
	// 	fmt.Println("No intent filters found.")
	// }

	fmt.Println("===== End of Report =====")
}
