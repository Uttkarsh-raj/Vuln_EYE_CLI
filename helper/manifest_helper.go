package helper

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Uttkarsh-raj/veye/models"
)

func AnalyzeManifest() error {
	file, err := os.Open("./app/src/main/AndroidManifest.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	manifest := &models.ManifestData{}

	permissionPattern := regexp.MustCompile(`android:name="(android\.permission\.[A-Z_]+)"`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "<uses-permission") {
			matches := permissionPattern.FindStringSubmatch(line)
			if len(matches) == 2 {
				permission := matches[1]
				manifest.Permissions = append(manifest.Permissions, permission)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	logPermissionsWithWarning(manifest)

	return nil
}

func logPermissionsWithWarning(manifest *models.ManifestData) {
	fmt.Println("===== Permissions Analysis =====")

	if len(manifest.Permissions) > 0 {
		fmt.Println("Permissions found in the manifest:")
		for _, permission := range manifest.Permissions {
			// Print permission
			fmt.Printf("Permission granted: %s\n", permission)

			explainPermissionRisk(permission)
		}
	} else {
		fmt.Println("No permissions found.")
	}

	fmt.Println("===== End of Permissions Analysis =====")
}

func explainPermissionRisk(permission string) {
	switch permission {
	case "android.permission.READ_CONTACTS":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows the app to read the user's contacts data, which is sensitive.\n")
		fmt.Println("  Mitigation: Only request this permission if absolutely necessary and explain to the user why it is needed.\033[0m")
	case "android.permission.ACCESS_FINE_LOCATION":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows the app to access precise location data, which can be sensitive.\n")
		fmt.Println("  Mitigation: Minimize location requests and ensure data is handled securely.\033[0m")
	case "android.permission.READ_EXTERNAL_STORAGE":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Grants access to the user's external storage, potentially exposing private files.\n")
		fmt.Println("  Mitigation: Limit the usage to required files only, and ensure proper permission checks are in place.\033[0m")
	case "android.permission.WRITE_EXTERNAL_STORAGE":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows the app to write to external storage, which can lead to data tampering or leaks.\n")
		fmt.Println("  Mitigation: Ensure proper access controls and avoid writing sensitive data to external storage.\033[0m")
	case "android.permission.CAMERA":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Grants access to the device's camera, which could be exploited for unauthorized photo or video capture.\n")
		fmt.Println("  Mitigation: Ensure users are aware when the camera is in use and avoid unnecessary camera access.\033[0m")
	case "android.permission.RECORD_AUDIO":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows the app to record audio, potentially leading to privacy concerns if misused.\n")
		fmt.Println("  Mitigation: Clearly inform users when recording is happening and request permission only when necessary.\033[0m")
	case "android.permission.SEND_SMS":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows the app to send SMS messages, which could result in spam or unauthorized messaging.\n")
		fmt.Println("  Mitigation: Use this permission sparingly and ensure users are informed about SMS sending activities.\033[0m")
	case "android.permission.READ_SMS":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Grants access to SMS messages, which could expose sensitive information.\n")
		fmt.Println("  Mitigation: Access SMS only when absolutely necessary and handle message data with care.\033[0m")
	case "android.permission.ACCESS_COARSE_LOCATION":
		fmt.Printf("Risk Explanation for %s:\n", permission)
		fmt.Printf("\033[33m  ⚠️  Risk: Allows access to approximate location data, which could still be used for tracking.\n")
		fmt.Println("  Mitigation: Use location data responsibly and only when required.\033[0m")
	default:
	}
	fmt.Println("--------------------------------------------")
}
