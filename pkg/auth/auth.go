package auth

import "fmt"

type Auth interface {
	SetDefaults()
	Validate() bool
	Authenticate() error
}

func Authenticate(auth Auth) error {
	auth.SetDefaults()

	if !auth.Validate() {
		return fmt.Errorf("Validation failed for %v", auth)
	}

	return auth.Authenticate()
}
