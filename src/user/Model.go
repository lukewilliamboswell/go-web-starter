package user

// User represents a user in the database
type User struct {

	// An identifier for the caller set by the identity provider.
	PrincipalId string `json:"X-MS-CLIENT-PRINCIPAL-IDP"`

	// A human-readable name for the caller set by the identity provider.
	PrincipalName string `json:"X-MS-CLIENT-PRINCIPAL-NAME"`

	// The name of the identity provider used by App Service Authentication.
	PrincipalProvider string `json:"X-Ms-Client-Principal-Idp"`

	// A Base64 encoded JSON representation of available claims.
	PrincipalClaims string `json:"X-Ms-Client-Principal"`
}
