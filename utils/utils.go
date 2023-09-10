package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"gopkg.in/gomail.v2"
)

func GetSha256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ValidateEmail(email string) bool {
	// Regular expression to match email format
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`

	// Compile the regex pattern
	pattern := regexp.MustCompile(regex)

	// Use the pattern to match the email
	return pattern.MatchString(email)
}

func GetUniqueCustomerId() string {
	uuid := uuid.New()
	return uuid.String()
}

func TriggerVerificationEmail(toEmailId string) (bool, error) {

	templatePath := "./templates/emailtemplate.html"
	templateContent, err := template.ParseFiles(templatePath)
	if err != nil {
		return false, err
	}

	// Create a buffer to store the template output
	var tplBuffer bytes.Buffer
	data := struct {
		Link string
	}{
		Link: secret.HOST + strings.Replace(types.VERIFY_EMAIL_PATH, secret.EMAIL_PATH_STRING, toEmailId, -1),
	}
	err = templateContent.Execute(&tplBuffer, data)
	if err != nil {
		return false, err
	}

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", secret.FROM_EMAIL)
	m.SetHeader("To", toEmailId)
	m.SetHeader("Subject", "Verify your account on TinyFundConnect")
	m.SetBody("text/html", tplBuffer.String())

	// Create a new dialer for the SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, secret.FROM_EMAIL, secret.FROM_EMAIL_PASSWORD)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}

func ValidatePassword(password string) bool {
	// Check for at least one capital letter
	hasCapital := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasCapital = true
			break
		}
	}

	// Check for at least one digit
	hasDigit := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigit = true
			break
		}
	}

	// Check for at least one special character
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/~"
	for _, char := range password {
		if containsChar(specialChars, char) {
			hasSpecial = true
			break
		}
	}

	// Check for length greater than 8
	isLongEnough := len(password) > 8

	return hasCapital && hasDigit && hasSpecial && isLongEnough
}

func containsChar(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}

func GetBytesFromInterface(in interface{}) []byte {
	jsonData, err := json.Marshal(in)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte{}
	}

	return jsonData

}

func GenericEncodeToMap(m interface{}) ([]map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface []map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
