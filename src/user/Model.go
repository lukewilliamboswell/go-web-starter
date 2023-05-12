package user

// User represents a user in the database
type User struct {

	// An identifier for the caller set by the identity provider.
	PrincipalId string `json:"X-Ms-Client-Principal-Id"`

	// A human-readable name for the caller set by the identity provider.
	PrincipalName string `json:"X-Ms-Client-Principal-Name"`

	// The name of the identity provider used by App Service Authentication.
	PrincipalProvider string `json:"X-Ms-Client-Principal-Idp"`

	// Application Access Level
	AccessLevel AuthCode `json:"access_level"`
}

type AuthCode int

const (
	DENY_ACCESS  AuthCode = 0
	USER_ACCESS  AuthCode = 1
	ADMIN_ACCESS AuthCode = 2
)
