package utilities

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

// IsValidName
func IsValidName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_ \-\'\\.]{1,100}$`)
	return len(name) >= 1 && len(name) <= 100 && re.MatchString(name) && !strings.Contains(name, "  ") && !strings.HasPrefix(name, " ")
}

// IsValidEmail
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return len(email) >= 5 && len(email) <= 100 && (err == nil)
}

// IsValidPassword
func IsValidPassword(password string) bool {
	return len(password) >= 6 && len(password) <= 25
}

// IsValidBirthDate
func IsValidBirthDate(birthdate string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(birthdate)
}

// IsValidDescription
func IsValidDescription(description string) bool {
	return len(description) <= 2048
}

func ToInt(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// check if size biger in kb
// check if request size is bigger than limit in KB
func IsSizeBiggerThan(size float64, r *http.Request, w http.ResponseWriter, errormsg string) bool {
	maxBytes := int64(size * 1024)

	if r.ContentLength > maxBytes {
		WriteJSON(w, http.StatusBadRequest, errormsg, nil)
		return true
	}

	return false
}
