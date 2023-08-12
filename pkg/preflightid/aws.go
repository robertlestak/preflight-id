package preflightid

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
)

type IDProviderAWS struct {
	ARN string `json:"arn" yaml:"arn"`
}

func (p *IDProviderAWS) Run() error {
	l := log.WithFields(log.Fields{
		"preflight": "id",
		"provider":  "aws",
	})
	l.Debug("running preflight-id")
	if p.ARN == "" {
		return errors.New("ARN not configured")
	}
	// Create a new AWS session using environment credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create an STS client
	svc := sts.New(sess)
	// Get caller identity
	resp, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	if *resp.Arn != p.ARN {
		failStr := fmt.Sprintf("failed - expected %s, got %s", p.ARN, *resp.Arn)
		l.Error(failStr)
		return errors.New(failStr)
	}
	l.Info("passed")
	return nil
}
