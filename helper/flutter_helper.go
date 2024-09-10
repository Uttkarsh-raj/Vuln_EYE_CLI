package helper

import (
	"regexp"
	"strings"
)

// helper function for the versions
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
