package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
)

type AssumeRole struct {
	duration        *int64
	roleArn         *string
	roleSessionName *string
	session         *SharedSession
	profile         string
	region          string
}

func NewAssumeRole(duration int64, roleArn string,
	roleSessionName string, session *SharedSession,
	profile string, region string) (*AssumeRole, error) {
	return &AssumeRole{
		duration:        &duration,
		roleArn:         &roleArn,
		roleSessionName: &roleSessionName,
		session:         session,
		profile:         profile,
		region:          region,
	}, nil
}

func (a *AssumeRole) AssumeRoleAndSaveProfile() error {
	input := &sts.AssumeRoleInput{
		DurationSeconds: a.duration,
		RoleArn:         a.roleArn,
		RoleSessionName: a.roleSessionName,
	}

	err := input.Validate()
	if err != nil {
		return fmt.Errorf("RoleInput not valid: %s", err)
	}

	output, err := a.session.svc.AssumeRole(input)
	if err != nil {
		return fmt.Errorf("Error on AssumeRoleInput: %s", err)
	}

	ini, err := NewConfigFile()
	if err != nil {
		return fmt.Errorf("Error on AWS config file: %s", err)
	}

	err = ini.EnsureSectionAndSave("profile "+a.profile, map[string]string{
		"output": "json",
		"region": a.region,
	})

	if err != nil {
		return err
	}

	ini, err = NewCredFile()
	if err != nil {
		return fmt.Errorf("Error on AWS credentials file: %s", err)
	}

	err = ini.EnsureSectionAndSave(a.profile, map[string]string{
		"aws_access_key_id":     *output.Credentials.AccessKeyId,
		"aws_secret_access_key": *output.Credentials.SecretAccessKey,
		"aws_session_token":     *output.Credentials.SessionToken,
	})

	if err != nil {
		return err
	}

	return nil
}
