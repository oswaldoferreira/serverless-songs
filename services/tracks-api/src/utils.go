package utils

import (
	"fmt"
)

type Header map[string]string

// JSONHeader foo
var JSONHeader Header = Header{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}

// GetUserID returns the user ID.
// One can add extra context at the custom authorizer
// to get more data from it (e.g. additional JWT token payload).
func GetUserID(auth map[string]interface{}) string {
	userID := fmt.Sprintf("%s", auth["principalId"])

	return userID
}
