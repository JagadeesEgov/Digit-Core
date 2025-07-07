package utils

import (
	"regexp"
)

// URLValidator provides URL validation functionality
type URLValidator struct{}

// NewURLValidator creates a new URLValidator instance
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateURL validates a URL using the same regex pattern as the Java implementation
func (v *URLValidator) ValidateURL(url string) bool {
	// Same regex pattern as the Java implementation
	urlRegex := `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`
	
	matched, err := regexp.MatchString(urlRegex, url)
	if err != nil {
		return false
	}
	
	return matched
}