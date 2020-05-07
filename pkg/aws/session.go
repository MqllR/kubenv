package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type SharedSession struct {
	AwsProfile string
	Sess       *session.Session
}

func NewSharedSession(awsProfile string) (*SharedSession, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})

	if err != nil {
		return nil, err
	}

	return &SharedSession{
		AwsProfile: awsProfile,
		Sess:       sess,
	}, nil
}

func (sess *SharedSession) IsExpired() bool {
	svc := sts.New(sess.Sess)
	input := &sts.GetCallerIdentityInput{}
	_, err := svc.GetCallerIdentity(input)

	if err != nil {
		return true
	}

	return false
}
