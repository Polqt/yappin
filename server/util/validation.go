package util

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail checks if the email is valid
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	
	// Allow alphanumeric characters and underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	
	return nil
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters long")
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes and control characters except newlines and tabs
	cleaned := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || (r >= 32 && r != 127) {
			return r
		}
		return -1
	}, input)
	
	return strings.TrimSpace(cleaned)
}

// ValidateUUID checks if the string is a valid UUID format
func ValidateUUID(uuid string) error {
	matched, _ := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", strings.ToLower(uuid))
	if !matched {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}
