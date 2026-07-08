package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func MakeSlugSimple(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "post"
	}
	uuidPart := uuid.New().String()[:8]
	return s + "-" + uuidPart
}
