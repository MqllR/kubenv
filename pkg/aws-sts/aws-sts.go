package awssts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/aws"
)

type AssumeRole struct {
	duration        *int64
	roleArn         *string
	roleSessionName *string
	session         *aws.SharedSession
	profile         string
	region          string
}

func NewAssumeRole(duration int64, roleArn string,
	roleSessionName string, session *aws.SharedSession,
	profile string, region string) *AssumeRole {
	return &AssumeRole{
		duration:        &duration,
		roleArn:         &roleArn,
		roleSessionName: &roleSessionName,
		session:         session,
		profile:         profile,
		region:          region,
	}
}

func (a *AssumeRole) Authenticate() error {
	input := &sts.AssumeRoleInput{
		DurationSeconds: a.duration,
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

	klog.V(5).Info("Token received: %v", output)

	ini, err := aws.NewConfigFile()
	if err != nil {
		return fmt.Errorf("Error on AWS config file: %s", err)
	}

	err = ini.EnsureSectionAndSave("profile "+a.profile, map[string]string{
		"output": "json",
		"region": a.region,
	})

	klog.V(2).Infof("Profile %s saved in AWS config file", a.profile)

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

	klog.V(2).Infof("Credentials %s saved in AWS credentials file", a.profile)

	return nil
}
