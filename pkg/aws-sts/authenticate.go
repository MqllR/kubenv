package awssts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/mqllr/kubenv/pkg/aws"
)

type AssumeRole struct {
	Duration        *int64
	roleArn         *string
	roleSessionName *string
	session         *aws.SharedSession
	profile         string
	region          string
}

func NewAssumeRole(roleArn string, roleSessionName string,
	session *aws.SharedSession, profile string) *AssumeRole {

	return &AssumeRole{
		roleArn:         &roleArn,
		roleSessionName: &roleSessionName,
		session:         session,
		profile:         profile,
	}
}

func (a *AssumeRole) SetAWSRegion(awsRegion string) {
	a.region = awsRegion
}

func (a *AssumeRole) SetDuration(duration *int64) {
	a.Duration = duration
}

func (a *AssumeRole) SetDefaults() {
	if *a.Duration == 0 {
		a.Duration = &DefaultDuration
	}

	if a.region == "" {
		a.region = DefaultAWSRegion
	}
}

func (a *AssumeRole) Validate() bool {
	return true
}

func (a *AssumeRole) Authenticate() error {
	input := &sts.AssumeRoleInput{
		DurationSeconds: a.Duration,
		RoleArn:         a.roleArn,
		RoleSessionName: a.roleSessionName,
	}

	err := input.Validate()
	if err != nil {
		return fmt.Errorf("RoleInput not valid: %s", err)
	}

	output, err := a.session.Svc.AssumeRole(input)
	if err != nil {
		return fmt.Errorf("Error on AssumeRoleInput: %s", err)
	}

	ini, err := aws.NewConfigFile()
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

	ini, err = aws.NewCredFile()
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
