package utils

import (
	"fmt"
	"regexp"
)

// Regular expression to parse gRPC service name and method name.
var re *regexp.Regexp = regexp.MustCompile(`/([a-zA-Z0-9\.]+)/([a-zA-Z0-9\.]+)(.*)`)

func ParseGrpcPath(path string) (servicename string, methodname string, err error) {
	matches := re.FindStringSubmatch(path)
	if len(matches) != 4 || matches[3] != "" {
		return "", "", fmt.Errorf("Failed to match path")
	}
	return matches[1], matches[2], nil
}
