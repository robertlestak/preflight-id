package preflightid

import (
	"errors"
	"fmt"
	"strings"

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
	if *resp.Arn == p.ARN {
		// exact match
		l.Info("passed")
		return nil
	}
	// check if we have an assumed role
	// the resp.Arn will be in the format of
	// arn:aws:sts::123456789012:assumed-role/role-name/role-session-name
	if strings.Contains(*resp.Arn, "assumed-role/") {
		// check if the ARN matches the assumed role ARN
		roleName := strings.Split(strings.Split(*resp.Arn, "/")[1], "/")[0]
		accountNumber := strings.Split(*resp.Arn, ":")[4]
		assumedRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountNumber, roleName)
		if assumedRoleARN == p.ARN {
			l.Info("passed")
			return nil
		}
	}
	failStr := fmt.Sprintf("failed - expected %s, got %s", p.ARN, *resp.Arn)
	l.Error(failStr)
	return errors.New(failStr)
}
