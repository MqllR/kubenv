package config

var (
	AvailableAuthProviders = []string{
		"aws-google-auth",
		"aws-azure-login",
		"aws-sts",
	}

	AvailableK8SSyncMode = []string{
		"local",
	}
)
