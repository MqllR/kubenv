package awssts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/mqllr/kubenv/pkg/aws"
)

var (
	DefaultDuration int64 = 3600
)

type AssumeRole struct {
	Duration        *int64
	roleArn         *string
	roleSessionName *string
	session         *aws.SharedSession
	profile         string
	region          string
}

func NewAssumeRole(roleArn string,
	roleSessionName string, session *aws.SharedSession,
	profile string, region string) *AssumeRole {
	return &AssumeRole{
		roleArn:         &roleArn,
		roleSessionName: &roleSessionName,
		session:         session,
		profile:         profile,
		region:          region,
	}
}

func (a *AssumeRole) SetDefaults() {
	if *a.Duration == 0 {
		a.Duration = &DefaultDuration
	}
}

func (a *AssumeRole) Authenticate() error {
	a.SetDefaults()

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
