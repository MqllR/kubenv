package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type SharedSession struct {
	Svc *sts.STS
}

func NewSharedSession(awsProfile string) (*SharedSession, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})

	if err != nil {
		return nil, err
	}

	return &SharedSession{
		Svc: sts.New(sess),
	}, nil
}

func (s *SharedSession) IsExpired() bool {
	input := &sts.GetCallerIdentityInput{}
	_, err := s.Svc.GetCallerIdentity(input)

	if err != nil {
		return true
	}

	return false
}
