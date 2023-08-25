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
	ARN   string `json:"arn" yaml:"arn"`
	Equiv bool   `json:"equiv" yaml:"equiv"`
}

func (p *IDProviderAWS) RunEquiv() bool {
	return p.Equiv
}

func (p *IDProviderAWS) Equivalent() {
	l := Logger
	l.Debug("printing equivalent command")
	cmd := `ID=$(aws sts get-caller-identity --query Arn --output text);`
	// if it contains "assumed-role/" then it's an assumed role and we need to parse it
	cmd += `if [[ $ID == *"assumed-role/"* ]]; then ROLE_NAME=$(echo $ID | cut -d/ -f2); ACCOUNT_NUMBER=$(echo $ID | cut -d: -f5); ARN="arn:aws:iam::$ACCOUNT_NUMBER:role/$ROLE_NAME"; else ARN=$ID; fi;`
	cmd += fmt.Sprintf(`if [ "$ARN" != "%s" ]; then echo "ARN $ARN does not match expected %s"; exit 1; fi`, p.ARN, p.ARN)
	cmd = fmt.Sprintf("sh -c '%s'", cmd)
	fmt.Println(cmd)
}

func (p *IDProviderAWS) Run() error {
	l := Logger.WithFields(log.Fields{
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
